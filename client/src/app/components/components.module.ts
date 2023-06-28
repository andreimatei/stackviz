import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule } from '@angular/forms';
import { StacksComponent } from "./stacks/stacks.component";

@NgModule({
  imports: [
    CommonModule,
    ReactiveFormsModule,
    StacksComponent,
  ],
  exports: [
    StacksComponent,
  ]
})
export class ComponentsModule {
}
