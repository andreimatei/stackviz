import { Component, EventEmitter, Output } from "@angular/core";
import { NestedTreeControl } from "@angular/cdk/tree";
import { GetTypeInfoGQL, TypeInfo, VarInfo } from "../../graphql/graphql-codegen-generated";
import { CollectionViewer, DataSource, SelectionChange } from "@angular/cdk/collections";
import { BehaviorSubject, map, merge, Observable } from "rxjs";
import { MatCheckboxChange, MatCheckboxModule } from "@angular/material/checkbox";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { NgIf } from "@angular/common";
import { MatIconModule } from "@angular/material/icon";
import { MatButtonModule } from "@angular/material/button";
import { MatTreeModule } from "@angular/material/tree";

@Component({
    selector: 'type-info',
    templateUrl: './type-info.component.html',
    styleUrls: ['type-info.component.css'],
    standalone: true,
    imports: [
        MatTreeModule,
        MatCheckboxModule,
        MatButtonModule,
        MatIconModule,
        NgIf,
        MatProgressBarModule,
    ],
})
export class TypeInfoComponent {
  dataSource: TypesDataSource;
  treeControl: NestedTreeControl<TreeNode>;

  _exprs: string[] = [];

  get exprs(): string[] {
    return this._exprs;
  }

  set exprs(val) {
    this.dataSource.exprs = val;
    this._exprs = val;
  }

  @Output()
  checkedChange = new EventEmitter<CheckedEventArg>();

  constructor(
    private readonly typeQuery: GetTypeInfoGQL,
  ) {
    this.treeControl = new NestedTreeControl<TreeNode>(node => node.children);
    this.dataSource = new TypesDataSource(this.treeControl, typeQuery);
  }

  hasChild = (_: number, node: TreeNode) => node.expandable;

  onCheckedChange(ev: MatCheckboxChange, node: TreeNode) {
    const cev = new CheckedEventArg(node.expr, ev.checked);
    this.checkedChange.emit(cev);
  }
}

export class CheckedEventArg {
  constructor(public expr: string, public checked: boolean) {
  }
}

export class TreeNode {
  children: TreeNode[];
  color: string;
  fontWeight: string;
  isLoading: boolean = false;

  constructor(
    readonly name: string,
    readonly expr: string,
    readonly type: string,
    public expandable: boolean,
    readonly checked: boolean,
    color?: string, fontWeight?: string,
  ) {
    this.children = [];
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
}

export class TypesDataSource implements DataSource<TreeNode> {
  // dataChange will be notified whenever there is a change that causes the tree
  // to be reloaded.
  dataChange = new BehaviorSubject<TreeNode[]>([]);
  vars: VarInfo[] = [];
  types: Map<string, TypeInfo> = new Map();
  exprs: string[] = [];

  constructor(
    private _treeControl: NestedTreeControl<TreeNode>,
    private _database: GetTypeInfoGQL,
  ) {
  }

  get data(): TreeNode[] {
    return this.dataChange.value;
  }

  set data(val: TreeNode[]) {
    this.dataChange.next(val);
  }

  updateData(data: TreeNode[]) {
    // Clear the data, otherwise the tree's change detection doesn't pick up
    // changes in nested objects.
    this.data = [];
    this.data = data;
  }

  initData(vars: VarInfo[], types: TypeInfo[], exprs: string[]) {
    this.vars = vars;
    for (const t of types) {
      this.types.set(t.Name, t);
    }
    this.exprs = exprs;
    const data = vars.map<TreeNode>(v => {
      const ti = types.find(t => t.Name == v.Type);
      const expandable = !ti || (ti.FieldsNotLoaded || ti.Fields!.length > 0);
      const checked = exprs.includes(v.Name);
      const n: TreeNode = new TreeNode(
        v.Name, v.Name /* expr */, v.Type,
        expandable,
        checked,
        v.LoclistAvailable ? 'black' : 'gray',
        v.FormalParameter ? 'bold' : 'normal',
      );
      return n
    })
    this.updateData(data);
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
        if (node.children.length > 0) {
          // children already populated; nothing more to do.
          return
        }

        const ti = this.types.get(node.type);
        console.log("expanding: ti:", node.type, ti)
        // if (!ti) {
        //   // We failed to find the type. Make the node a leaf so that we don't
        //   // attempt to load it again.
        //   node.expandable = false;
        //   return;
        // }
        if (!ti || ti.FieldsNotLoaded) {
          console.log("loading children for %s...", node.type);
          node.isLoading = true;
          this._database.fetch({name: node.type}).pipe(
            map(res => {
              console.log("loading children... done", res)
              if (res.error) {
                console.log(res.error);
                return;
              }
              this.types.set(node.type, res.data.typeInfo);
              if (!res.data.typeInfo.Fields) {
                return;
              }
              node.children = this.typeInfoToTreeNodes(res.data.typeInfo, node.expr, this.exprs);
              node.expandable = node.children.length > 0;
              console.log("new children: ", node.children)
            })
          ).subscribe(_ => {
            // notify the change
            this.updateData(this.data);
            node.isLoading = false;
          })
        } else {
          node.children = this.typeInfoToTreeNodes(ti, node.expr, this.exprs);
          node.expandable = node.children.length > 0;
          // notify the change
          this.updateData(this.data);
        }
      });
    }
  }

  // !!!
  // getType(typename: string): TypeInfo {
  //   const ti = this.types.get(node.type);
  //   if (ti) {
  //     return ti;
  //   }
  // }

  typeInfoToTreeNodes(ti: TypeInfo, path: string, exprs: string[]): TreeNode[] {
    if (ti.FieldsNotLoaded) {
      console.log("unexpected fields not loaded")
      return []
    }
    const res: TreeNode[] = [];
    for (const f of ti.Fields!) {
      if (!f) continue;
      let expr = path + "." + f.Name;
      const cti = this.types.get(f.Type)
      const expandable = !cti || (cti.FieldsNotLoaded || cti.Fields!.length > 0);
      const checked = exprs.includes(expr);
      const n: TreeNode = new TreeNode(f.Name, expr, f.Type, expandable, checked);
      res.push(n);
    }
    return res;
  }
}
