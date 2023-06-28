import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { StacksComponent } from "./stacks/stacks.component";
import { DataTableModule } from 'traceviz/dist/ngx-traceviz-lib';


@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,
    StacksComponent,
    DataTableModule,
  ],
  exports: [
    StacksComponent,
  ]
})
export class ComponentsModule {
}
