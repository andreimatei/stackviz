import { Component, EventEmitter, Inject, Input, Output } from "@angular/core";
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
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogModule,
  MatDialogRef
} from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { FormsModule } from "@angular/forms";
import { MatRadioModule } from "@angular/material/radio";

// FlightRecorderEventSpec represents the specification used for recording data.
// It is part of a FrameSpec.
export interface FlightRecorderEventSpec {
  Expr: string;
  KeyExpr: string;
}

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
    MatDialogModule,
  ],
})
export class TypeInfoComponent {
  dataSource: TypesDataSource;
  treeControl: NestedTreeControl<TreeNode>;

  @Input()
  public flightRecorderEventSpecs!: FlightRecorderEventSpec[];

  // !!!
  // _exprs: string[] = [];
  //
  // get exprs(): string[] {
  //   return this._exprs;
  // }
  //
  // set exprs(val) {
  //   this.dataSource.exprs = val;
  //   this._exprs = val;
  // }

  @Output()
  checkedChange = new EventEmitter<CheckedEvent>();
  @Output()
  flightRecorderChange = new EventEmitter<FlightRecorderEvent>();

  constructor(
    private readonly typeQuery: GetTypeInfoGQL,
    public dialog: MatDialog,
  ) {
    this.treeControl = new NestedTreeControl<TreeNode>(node => node.children);
    this.dataSource = new TypesDataSource(this.treeControl, typeQuery);
  }

  hasChild = (_: number, node: TreeNode) => node.expandable;

  onCheckedChange(ev: MatCheckboxChange, node: TreeNode) {
    const cev = new CheckedEvent(node.expr, ev.checked);
    this.checkedChange.emit(cev);
  }

  openFlightRecorderDialog(varName: any): void {
    console.log("opening dialog for var:", varName);
    const dialogRef = this.dialog.open(
      FlightRecorderDialog, {
        data: {
          varName: varName,
          // TODO(andrei): Can there be multiple event specs for the same
          // variable (with different key expressions)? I guess the UI doesn't
          // allow creating them.
          existingSpec: this.flightRecorderEventSpecs.find(
            s => s.expr == varName),
        },
      });
    dialogRef.afterClosed().subscribe(result => {
      if (result === undefined) {
        return;
      }
      const ev = result as FlightRecorderEvent;
      console.log('The dialog was closed with result:', ev);
      this.flightRecorderChange.emit(ev);
    });
  }
}

export class CheckedEvent {
  constructor(public expr: string, public checked: boolean) {
  }
}

export const goroutineIDKey = Symbol('goroutineID');

export class FlightRecorderEvent {
  constructor(public expr: string, public key: string | typeof goroutineIDKey, public deleted: boolean) {
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
    // formalParam indicates whether this node corresponds to a formal parameter
    // of the function. Parameters are shown in bold and get a flight-recorder
    // button.
    readonly formalParam: boolean,
    // formalParamRecursive is set for formal parameters and their children,
    // recursively. Nodes with this flag set get the flight-recorder button.
    readonly formalParamRecursive: boolean,
    readonly loclistAvailable: boolean,
  ) {
    console.log("!!! TreeNode:", name, expr, type, expandable, checked, formalParam, loclistAvailable);
    this.color = loclistAvailable ? 'black' : 'gray';
    this.fontWeight = formalParam ? 'bold' : 'normal';
    this.children = [];
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
    console.log("!!! initData:", exprs);
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
        v.FormalParameter,
        v.FormalParameter,
        v.LoclistAvailable,
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
              node.children = this.typeInfoToTreeNodes(node, res.data.typeInfo, node.expr, this.exprs);
              node.expandable = node.children.length > 0;
              console.log("new children: ", node.children)
            })
          ).subscribe(_ => {
            // notify the change
            this.updateData(this.data);
            node.isLoading = false;
          })
        } else {
          node.children = this.typeInfoToTreeNodes(node, ti, node.expr, this.exprs);
          node.expandable = node.children.length > 0;
          // notify the change
          this.updateData(this.data);
        }
      });
    }
  }

  typeInfoToTreeNodes(
    parent: TreeNode, ti: TypeInfo, path: string, exprs: string[],
  ): TreeNode[] {
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
      const n: TreeNode = new TreeNode(
        f.Name, expr, f.Type, expandable, checked,
        false, // formalParam
        parent.formalParamRecursive,
        parent.loclistAvailable // if the parent is available, so is its field
      );
      res.push(n);
    }
    return res;
  }
}

@Component({
  selector: 'flight-recorder-dialog',
  templateUrl: 'flight-recorder-dialog.component.html',
  styleUrls: ['flight-recorder-dialog.component.css'],
  standalone: true,
  imports: [MatDialogModule, MatFormFieldModule, MatInputModule, FormsModule, MatButtonModule, MatRadioModule],
})
export class FlightRecorderDialog {
  keyOption: string = "goroutineID";
  customKeyExpr: string = "";

  constructor(
    public dialogRef: MatDialogRef<FlightRecorderDialog>,
    @Inject(MAT_DIALOG_DATA) public data: {
      varName: string,
      existingSpec?: FlightRecorderEventSpec,
    },
  ) {
    console.log("creating dialog with existing spec:", data);
  }

  onCancel(): void {
    this.dialogRef.close();
  }

  onOk(): void {
    console.log("closing dialog for var:", this.data.varName, this);
    this.dialogRef.close(new FlightRecorderEvent(
      this.data.varName,
      this.keyOption == 'goroutineID' ? goroutineIDKey : this.customKeyExpr,
      // TODO(andrei): add delete button
      false,
    ));
  }
}
