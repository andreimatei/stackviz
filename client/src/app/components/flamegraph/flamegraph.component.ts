import {
  AfterViewInit,
  Component, ElementRef,
  Input,
  OnInit,
  ViewChild,
  ViewEncapsulation
} from "@angular/core";
import * as d3 from 'd3';
import { FlameGraph, flamegraph } from './flamegraph-lib/index';
import { GetTreeGQL } from "../../graphql/graphql-codegen-generated";
import { tap } from "rxjs";

@Component({
  selector: 'app-flamegraph',
  standalone: true,
  imports: [],
  template: `
Flamegraph goes here
<div id="chart">
</div>
<div #details id="details">
</div>
  `,
  styleUrls: ['flamegraph.component.css'],
  // I've disabled encapsulation because otherwise the classes in the css file are not
  // available to the flamegraph svg. TODO(andrei): look into the proper solution.
  encapsulation: ViewEncapsulation.None,
})
export class FlamegraphComponent implements OnInit, AfterViewInit {
  @Input({required: false}) colID!: number;
  @Input({required: false}) snapID!: number;
  @ViewChild('details') details!: ElementRef;
  private flameGraph: FlameGraph;

  constructor(private getTreeQuery: GetTreeGQL) {
    this.flameGraph = flamegraph()
      .width(960)
      .cellHeight(18)
      .transitionDuration(750)
      // .minFrameSize(5)  // !!! I've removed this in order to make it look like traceviz
      // !!! .transitionEase(d3.easeCubic)
      .inverted(true)
      .sort(true)
      .title("")
      .onClick(this.onClick)
      .onCtrlClick(this.onCtrlClick)
      // !!! .differential(false)
      .selfValue(false);
  }

  ngAfterViewInit(): void {
    this.flameGraph.setDetailsElement(this.details.nativeElement);
  }

  ngOnInit(): void {
    this.redraw();
  }

  redraw() {
    this.getTreeQuery.fetch({colID: this.colID!, snapID: this.snapID!})
      .subscribe(value => {
        const data = JSON.parse(value.data.getTree);
        d3.select("#chart")
          .datum(data)
          .call(this.flameGraph);
      })

    // d3.json("assets/stacks.json").then(
    //   data => {
    //     d3.select("#chart")
    //       .datum(data)
    //       .call(flameGraph);
    //   }
    // )
    // //   , function(error: any, data: any) {
    // //   if (error) return console.warn(error);
    // //   d3.select("#chart")
    // //     .datum(data)
    // //     .call(flameGraph);
    // // });

  }

  onClick(node: any) {
    console.log("click on ", node);
  }

  onCtrlClick(node: any, ev: MouseEvent) {
    console.log("ctrl-click on ", node);
  }
}
