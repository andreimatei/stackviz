import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import {
  AddExprToCollectSpecGQL,
  GetAvailableVariablesGQL,
  GetCollectionGQL,
  GetTypeInfoGQL,
  ProcessSnapshot,
  RemoveExprFromCollectSpecGQL,
  TypeInfo,
  VarInfo
} from "../../graphql/graphql-codegen-generated";
import { ActivatedRoute } from "@angular/router";
import { AppCoreService, WeightedTreeComponent } from 'traceviz/dist/ngx-traceviz-lib';
import { Action, IntegerValue, Tree, Update, ValueMap, } from "traceviz-client-core";
import { MatDrawer } from "@angular/material/sidenav";
import { NestedTreeControl } from "@angular/cdk/tree";
import { MatCheckboxChange } from "@angular/material/checkbox";
import { ResizeEvent } from 'angular-resizable-element';
import { BehaviorSubject, map, merge, Observable } from "rxjs";
import { CollectionViewer, DataSource, SelectionChange } from "@angular/cdk/collections";


class TreeNode {
  children: TreeNode[];
  childrenNotLoaded: boolean = false;
  color: string;
  fontWeight: string;
  isLoading: boolean = false; // !!! needed?

  constructor(
    readonly name: string,
    readonly expr: string,
    readonly type: string,
    children: TreeNode[] | null,
    readonly checked: boolean,
    color?: string, fontWeight?: string,
  ) {
    if (children == null) {
      this.children = [];
    } else {
      this.children = children;
    }
    if (typeof color !== 'undefined') {
      this.color = color;
    } else {
      this.color = "black";
    }
    if (typeof fontWeight !== 'undefined') {
      this.fontWeight = fontWeight;
    } else {
      this.fontWeight = "normal";
    }
  }

  Size(): number {
    let n: number = 1;
    for (var c of this.children) {
      n += c.Size();
    }
    return n;
  }
}

class Frame {
  constructor(public name: string, public file: string, public line: number) {
  };
}

@Component({
  selector: 'snapshot',
  templateUrl: './snapshot.component.html',
  styleUrls: ['snapshot.component.css'],
  // Provide the AppCoreService at the level of this component, overriding the
  // `providedIn: 'root'` specified on the declaration of AppCoreService (which
  // asks for a single instance of AppCoreService to be injected everywhere in
  // the app). By scoping an instance to this SnapshotComponent, we prevent the
  // state from escaping this component. In turn, this means that every time
  // that a user navigates to the snapshot page, she gets an empty state (i.e.
  // any filtering they might have previously done on the page is gone; this is
  // considered a good thing.
  providers: [AppCoreService]
})
export class SnapshotComponent implements OnInit, AfterViewInit {
  protected collectionID!: number;
  protected snapshotID!: number;
  protected collectionName?: string;
  protected snapshots?: ProcessSnapshot[];
  @ViewChild(WeightedTreeComponent) weightedTree: WeightedTreeComponent | undefined;
  @ViewChild('functionDrawer') frameDetailsSidebar!: MatDrawer;
  dataSource: TypesDataSource;
  // !!! dataSource = new MatTreeNestedDataSource<TreeNode>();
  treeControl: NestedTreeControl<TreeNode>;

  // !!! unused?
  fetchTypeChildren(node: TreeNode): Observable<TreeNode[]> {
    this.treeControl.expansionModel.changed.subscribe()
    return this.typeQuery.fetch({name: node.name}).pipe(
      map(res => {
        if (res.error) {
          console.log(res.error)
          return [] as TreeNode[];
        }
        if (!res.data.typeInfo.Fields) {
          return [] as TreeNode[];
        }
        return typeInfoToTreeNodes(res.data.typeInfo, "" /* path */, [] /* exprs */);
      })
    )
  }

  protected selectedFrame?: Frame;
  // Data about the selected node. Each element is a string containing all the
  // captured variables from one frame (where all frames correspond to the
  // selected node).
  protected funcInfo?: string[];

  protected loadingAvailableVars: boolean = false;

  constructor(
    private readonly appCoreService: AppCoreService,
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly varsQuery: GetAvailableVariablesGQL,
    private readonly typeQuery: GetTypeInfoGQL,
    private readonly addExpr: AddExprToCollectSpecGQL,
    private readonly removeExpr: RemoveExprFromCollectSpecGQL,
    private readonly route: ActivatedRoute,
  ) {
    this.treeControl = new NestedTreeControl<TreeNode>(node => node.children);
    this.dataSource = new TypesDataSource(this.treeControl, typeQuery);

    // !!!
    const sumObserver = {
      sum: 0,
      next(value: any) {
        console.log('Next');
      },
      error() {
        console.log('error');
      },
      complete() {
        console.log('completed');
      }
    };
    this.typeQuery.fetch({name: "time.Time"}).subscribe(sumObserver)

  }

  ngOnInit(): void {
    // Get the collection ID and snapshot ID from the URL. The names of the URL
    // params are defined in the Routes collection.
    this.collectionID = Number(this.route.snapshot.paramMap.get('colID'));
    this.snapshotID = Number(this.route.snapshot.paramMap.get('snapID'));
    this.getCollectionQuery.fetch({colID: this.collectionID.toString()})
      .subscribe(results => {
        this.collectionName = results.data.collectionByID?.name;
        this.snapshots = results.data.collectionByID?.processSnapshots?.map(
          value => ({...value, snapshot: "",})
        )
      })
    this.appCoreService.appCore.globalState.set(
      "snapshot_id", new IntegerValue(this.snapshotID));
  }

  ngAfterViewInit() {
    this.weightedTree!.interactionsDir!.get().withAction(
      new Action(WeightedTreeComponent.NODE, WeightedTreeComponent.CTRL_CLICK,
        new Call(this.onNodeCtrlClick.bind(this))));
  }

  onNodeCtrlClick(localState: ValueMap): void {
    if (localState.has('vars')) {
      this.funcInfo = localState.expectStringList('vars');
    }

    const funcName = localState.expectString('full_name');
    this.selectedFrame = new Frame(funcName, localState.expectString('file'), localState.expectNumber('line'))
    const pcOffset = localState.expectNumber('pc_off');
    console.log("querying for available vars for func: %s off: %d", funcName, pcOffset);
    this.loadingAvailableVars = true;
    this.dataSource.data = [];
    this.varsQuery.fetch({func: funcName, pcOff: pcOffset})
      .subscribe(
        results => {
          console.log("got results");
          this.loadingAvailableVars = false;
          if (results.error) {
            console.log("!!! err: ", results.error)
            return
          }
          let exprs: string[] = results.data.frameInfo ? results.data.frameInfo.exprs : [];
          this.dataSource.data = convertToTree(
            results.data.availableVars.Vars,
            results.data.availableVars.Types,
            exprs);
          let numNodes: number = 0;
          for (var n of this.dataSource.data) {
            numNodes += n.Size()
          }
          console.log("tree size: %d", numNodes)
        })
    console.log("opening sidebar")
    this.frameDetailsSidebar.open();
  }

  closeSidebar() {
    this.frameDetailsSidebar.close()
  }

  onSelectedSnapshotChange(newValue: string) {
    let newSnapshotID = Number(newValue);
    this.appCoreService.appCore.globalState.get("snapshot_id").fold(new IntegerValue(newSnapshotID), false /* toggle */);
  }

  hasChild = (_: number, node: TreeNode) => node.childrenNotLoaded || node.children.length > 0;

  varChange(ev: MatCheckboxChange, node: TreeNode) {
    console.log(`clicked on ${node.name} in func ${this.selectedFrame}`);
    if (ev.checked) {
      this.addExpr.mutate({frame: this.selectedFrame!.name, expr: node.expr}).subscribe({
        next: value => console.log(value.data?.addExprToCollectSpec.frames![0].exprs)
      })
    } else {
      this.removeExpr.mutate({frame: this.selectedFrame!.name, expr: node.expr}).subscribe({
        next: value => console.log(value.data?.removeExprFromCollectSpec.frames![0].exprs)
      })
    }
  }

  public style: object = {
    "width": '500px',
    "text-color": 'red',
  };

  onResizeEnd(event: ResizeEvent): void {
    console.log("!!! width: %d", event.rectangle.width)
    this.style = {
      // position: 'fixed',
      // left: `${event.rectangle.left}px`,
      // top: `${event.rectangle.top}px`,
      // width: `${event.rectangle.width}px`,
      // height: `${event.rectangle.height}px`,
      width: `${event.rectangle.width}px`,
    };

  }

  validateResize(event: ResizeEvent): boolean {
    return true;
  }
}

// Call is an implementation of Update that calls the provided function.
class Call extends Update {
  constructor(private readonly handler: ((localState: ValueMap) => void)) {
    super();
  }

  override update(localState?: ValueMap | undefined) {
    this.handler(localState!)
  }

  get autoDocument(): string {
    throw new Error('Method not implemented.');
  }
}

function structToTree(ti: TypeInfo, types: TypeInfo[], path: string, level: number, exprs: string[]): TreeNode[] {
  // !!! protect against infinite recursion. I should do this based on a path.
  if (level == 2) {
    return [];
  }
  const res: TreeNode[] = [];
  for (var f of ti.Fields!) {
    if (!f) continue;
    let expr = path + "." + f.Name
    const n: TreeNode = new TreeNode(f.Name, expr, f.Type, null, exprs.includes(expr));
    const ti = types.find(t => t.Name == f!.Type);
    if (ti) {
      if (ti.FieldsNotLoaded) {
        console.log("!!! childrenNotLoaded on %s", f.Name)
        n.childrenNotLoaded = true
      } else {
        n.children = structToTree(ti, types, expr, level + 1, exprs)
      }
    }
    res.push(n);
  }
  return res;
}

function convertToTree(vars: VarInfo[], types: TypeInfo[], exprs: string[]): Array<TreeNode> {
  return vars.map<TreeNode>(v => {
    const n: TreeNode = new TreeNode(v.Name, v.Type, v.Name, null, exprs.includes(v.Name),
      v.LoclistAvailable ? 'black' : 'gray',
      v.FormalParameter ? 'bold' : 'normal',
    );
    const ti = types.find(t => t.Name == v.Type);
    if (ti) {
      if (ti.FieldsNotLoaded) {
        console.log("!!! childrenNotLoaded on var %s", v.Name)
        n.childrenNotLoaded = true
      } else {
        n.children = structToTree(ti, types, v.Name, 0, exprs)
      }
    }
    return n
  })
}

function typeInfoToTreeNodes(ti: TypeInfo, path: string, exprs: string[]): TreeNode[] {
  if (ti.FieldsNotLoaded) {
    console.log("unexpected fields not loaded")
    return []
  }
  const res: TreeNode[] = [];
  for (var f of ti.Fields!) {
    if (!f) continue;
    let expr = path + "." + f.Name;
    const n: TreeNode = new TreeNode(f.Name, f.Type, expr, null, exprs.includes(expr));
    n.childrenNotLoaded = true;
    res.push(n);
  }
  return res;
}

export class TypesDataSource implements DataSource<TreeNode> {
  // dataChange will be notified whenever there is a change that causes the tree
  // to be reloaded.
  dataChange = new BehaviorSubject<TreeNode[]>([]);

  get data(): TreeNode[] {
    return this.dataChange.value;
  }

  set data(val: TreeNode[]) {
    this.dataChange.next(val);
  }

  updateData(data: TreeNode[]) {
    this.data = [];
    this.data = data;
  }

  constructor(
    private _treeControl: NestedTreeControl<TreeNode>,
    private _database: GetTypeInfoGQL,
  ) {
  }

  connect(collectionViewer: CollectionViewer): Observable<TreeNode[]> {
    this._treeControl.expansionModel.changed.subscribe(change => {
      if (
        (change as SelectionChange<TreeNode>).added ||
        (change as SelectionChange<TreeNode>).removed
      ) {
        this.handleTreeControl(change as SelectionChange<TreeNode>);
      }
    });

    return merge(collectionViewer.viewChange, this.dataChange).pipe(map(() => this.data));
  }

  disconnect(collectionViewer: CollectionViewer): void {
  }

  /** Handle expand/collapse behaviors */
  handleTreeControl(change: SelectionChange<TreeNode>) {
    if (change.added) {
      change.added.forEach(node => {
        if (node.childrenNotLoaded) {
          console.log("loading children for %s...", node.type);

          node.isLoading = true;
          this._database.fetch({name: node.type}).pipe(
            map(res => {
              console.log("loading children... done", res)
              if (res.error) {
                console.log(res.error)
                return;
              }
              if (!res.data.typeInfo.Fields) {
                return;
              }
              node.children = typeInfoToTreeNodes(res.data.typeInfo, "" /* path */, [] /* exprs */);
              node.childrenNotLoaded = false;
              console.log("new children: ", node.children)
            })
          ).subscribe(_ => {
            // notify the change
            this.updateData(this.data);
            node.isLoading = false;
          })
        }
      });
    }
  }
}
