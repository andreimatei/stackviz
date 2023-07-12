import {
  AfterViewInit,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  Output,
  ViewChild,
  ViewEncapsulation
} from "@angular/core";
import * as d3 from 'd3';
import { FlameGraph, flamegraph } from './flamegraph-lib';
import { BehaviorSubject, merge, Observable, Subject } from "rxjs";
import { AngularResizeEventModule, ResizedEvent } from 'angular-resize-event';
import { Link } from "../../graphql/graphql-codegen-generated";

export interface Frame {
  name: string;
  details: string;
  file: string;
  line: number;
  pcoff: number;
  vars: Map<number, VarInfo[]>;
}

export interface VarInfo {
  Expr: string;
  Value: string;
  Links: Link[];
}

@Component({
  selector: 'app-flamegraph',
  standalone: true,
  imports: [AngularResizeEventModule],
  template: `
    <div #flamegraph id="flamegraph" class="flamegraph-container" (resized)="onResized($event)">
    </div>
    <div #details>
    </div>
  `,
  styleUrls: ['flamegraph.component.css'],
  // I've disabled encapsulation because otherwise the classes in the css file are not
  // available to the flamegraph svg. TODO(andrei): look into the proper solution.
  encapsulation: ViewEncapsulation.None,
})
export class FlamegraphComponent implements AfterViewInit {
  private _data$ = new BehaviorSubject<any>(null);
  private _resize$ = new Subject<{}>();

  @Input() set data$(val$: Observable<any>) {
    val$.subscribe(value => this._data$.next(value));
    merge(this._data$, this._resize$.asObservable()).subscribe(
      () => this.redraw(this._data$.getValue())
    );
  }

  @Output() ctrlClick: EventEmitter<Frame> = new EventEmitter<any>();

  @ViewChild('details') details!: ElementRef;
  @ViewChild('flamegraph') flameGraphElement!: ElementRef;
  private readonly flameGraph: FlameGraph;

  constructor() { // !!! private getTreeQuery: GetTreeGQL) {
    this.flameGraph = flamegraph()
      .width(1500)
      .cellHeight(18)
      .transitionDuration(750)
      // .minFrameSize(5)
      // !!! .transitionEase(d3.easeCubic)
      .inverted(true)
      .sort(true)
      .title("")
      .onClick(this.onClick)
      .onCtrlClick(this.onCtrlClick.bind(this))
      // !!! .differential(false)
      .selfValue(false);  // a node's value includes the values of its children
    console.log("!!! height", this.flameGraph.height());
  }

  ngAfterViewInit(): void {
    this.flameGraph.setDetailsElement(this.details.nativeElement);
  }

  onResized(event: ResizedEvent) {
    this.flameGraph.width(event.newRect.width);
    this._resize$.next({});
  }

  redraw(data: any): void {
    // HACK: reset the flamegraph's height, otherwise it stays the same height
    // as it was previously. This is because, if a size hasn't be explicitly
    // set, the flamegraph remembers its computed size after the first draw.
    this.flameGraph.height(null);
    // !!!clear the flamegraph if we receive null?
    if (data == null) {
      return;
    }
    d3.select("#flamegraph")
      .datum(data)
      .call(this.flameGraph);
  }

  onClick(node: any) {
  }

  onCtrlClick(node: any, ev: MouseEvent) {
    this.ctrlClick.emit(node.data as Frame);
  }
}
