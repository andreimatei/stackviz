import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TextboxComponent } from "./textbox/textbox.component";
import { ReactiveFormsModule } from '@angular/forms';

export { TextboxComponent } from "./textbox/textbox.component"

@NgModule({
  declarations: [
    TextboxComponent
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
  ],
  exports: [
    TextboxComponent
  ]
})
export class ComponentsModule { }
