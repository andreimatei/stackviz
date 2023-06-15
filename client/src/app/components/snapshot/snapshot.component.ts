import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import {
  AddExprToCollectSpecGQL, FrameInfo,
  GetAvailableVariablesGQL,
  GetAvailableVariablesQuery,
  GetCollectionGQL,
  ProcessSnapshot,
  RemoveExprFromCollectSpecGQL,
  TypeInfo
} from "../../graphql/graphql-codegen-generated";
import { ActivatedRoute } from "@angular/router";
import { AppCoreService, WeightedTreeComponent } from 'traceviz/dist/ngx-traceviz-lib';
import { Action, IntegerValue, Tree, Update, ValueMap, } from "traceviz-client-core";
import { MatDrawer } from "@angular/material/sidenav";
import { MatTreeNestedDataSource } from "@angular/material/tree";
import { NestedTreeControl } from "@angular/cdk/tree";
import { MatCheckboxChange } from "@angular/material/checkbox";
import { ResizeEvent } from 'angular-resizable-element';


interface TreeNode {
  name: string;
  expr: string;
  type: string;
  children: TreeNode[];
  checked: boolean;
}

class Frame {
  constructor(public name: string, public file: string, public line: number){};
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
  protected availableVars?: GetAvailableVariablesQuery['availableVars']['Vars'];
  @ViewChild(WeightedTreeComponent) weightedTree: WeightedTreeComponent | undefined;
  @ViewChild('functionDrawer') frameDetailsSidebar!: MatDrawer;
  dataSource = new MatTreeNestedDataSource<TreeNode>();
  treeControl = new NestedTreeControl<TreeNode>(node => node.children);


  protected selectedFrame?: Frame;
  // Data about the selected node. Each element is a string containing all the
  // captured variables from one frame (where all frames correspond to the
  // selected node).
  protected funcInfo?: string[];

  constructor(
    private readonly appCoreService: AppCoreService,
    private readonly getCollectionQuery: GetCollectionGQL,
    private readonly varsQuery: GetAvailableVariablesGQL,
    private readonly addExpr: AddExprToCollectSpecGQL,
    private readonly removeExpr: RemoveExprFromCollectSpecGQL,
    private readonly route: ActivatedRoute) {
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
    console.log(localState);
    console.log(this.frameDetailsSidebar!);
    if (localState.has('vars')) {
      console.log("!!! vars: ", localState.get('vars').toString());
      this.funcInfo = localState.expectStringList('vars');
    }

    const funcName = localState.expectString('full_name');
    this.selectedFrame = new Frame(funcName, localState.expectString('file'), localState.expectNumber('line'))
    const pcOffset = localState.expectNumber('pc_off');
    console.log("querying for available vars for func: %s off: %d", funcName, pcOffset);
    this.varsQuery.fetch({func: funcName, pcOff: pcOffset})
      .subscribe(
        results => {
          if (results.error) {
            console.log("!!! err: ", results.error)
            return
          }
          this.availableVars = results.data.availableVars.Vars;

          function structToTree(ti: TypeInfo, types: TypeInfo[], path: string, level: number, exprs: string[]): TreeNode[] {
            if (level == 3) {
              return [];
            }
            const res: TreeNode[] = [];
            for (var f of ti.Fields!) {
              if (!f) continue;
              let expr = path + "." + f.Name
              const n: TreeNode = {name: f.Name, type: f.Type, expr: expr, children: [], checked: exprs.includes(expr)};
              const ti = types.find(t => t.Name == f!.Type);
              if (ti) {
                n.children = structToTree(ti, types, expr, level + 1, exprs)
              }
              res.push(n);
            }
            return res;
          }

          function convertToTree(vars: Array<{
            Name: string;
            Type: string;
            VarType: number;
          }>, types: TypeInfo[], exprs: string[]): Array<TreeNode> {
            return vars.map<TreeNode>(v => {
              const n: TreeNode = {name: v.Name, type: v.Type, expr: v.Name, children: Array<TreeNode>(), checked: exprs.includes(v.Name)}
              const ti = types.find(t => t.Name == v.Type);
              if (ti) {
                n.children = structToTree(ti, types, v.Name, 0, exprs)
              }
              return n
            })
          }

          let exprs: string[] = results.data.frameInfo ? results.data.frameInfo.exprs : [];
          this.dataSource.data = convertToTree(
            results.data.availableVars.Vars,
            results.data.availableVars.Types,
            exprs);
        })

    this.frameDetailsSidebar.open();
  }

  closeSidebar() {
    this.frameDetailsSidebar.close()
  }

  onSelectedSnapshotChange(newValue: string) {
    let newSnapshotID = Number(newValue);
    this.appCoreService.appCore.globalState.get("snapshot_id").fold(new IntegerValue(newSnapshotID), false /* toggle */);
  }

  hasChild = (_: number, node: TreeNode) => node.children.length > 0;

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
