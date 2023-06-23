import { Component, EventEmitter, Output } from "@angular/core";
import { NestedTreeControl } from "@angular/cdk/tree";
import { GetTypeInfoGQL, TypeInfo } from "../../graphql/graphql-codegen-generated";
import { CollectionViewer, DataSource, SelectionChange } from "@angular/cdk/collections";
import { BehaviorSubject, map, merge, Observable } from "rxjs";
import { MatCheckboxChange } from "@angular/material/checkbox";

@Component({
  selector: 'type-info',
  templateUrl: './type-info.component.html',
  styleUrls: ['type-info.component.css'],
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

  hasChild = (_: number, node: TreeNode) => node.childrenNotLoaded || node.children.length > 0;

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

export class TypesDataSource implements DataSource<TreeNode> {
  // dataChange will be notified whenever there is a change that causes the tree
  // to be reloaded.
  dataChange = new BehaviorSubject<TreeNode[]>([]);
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
              node.children = this.typeInfoToTreeNodes(res.data.typeInfo, node.expr, this.exprs);
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

  typeInfoToTreeNodes(ti: TypeInfo, path: string, exprs: string[]): TreeNode[] {
    if (ti.FieldsNotLoaded) {
      console.log("unexpected fields not loaded")
      return []
    }
    const res: TreeNode[] = [];
    for (var f of ti.Fields!) {
      if (!f) continue;
      let expr = path + "." + f.Name;
      const n: TreeNode = new TreeNode(f.Name, expr, f.Type, null, exprs.includes(expr));
      n.childrenNotLoaded = true;
      res.push(n);
    }
    return res;
  }
}
