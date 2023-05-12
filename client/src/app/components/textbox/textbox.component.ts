import {Component, ContentChild} from '@angular/core';
import { FormControl } from '@angular/forms';

import { InteractionsDirective } from 'traceviz/dist/ngx-traceviz-lib';
import { Interactions, ValueMap, StringValue } from 'traceviz-client-core';

const TEXTBOX = 'textbox';
const TEXT_CHANGED = 'text_changed';

const supportedActions = new Array<[string, string]>(
    [TEXTBOX, TEXT_CHANGED],
);

@Component({
  selector: 'textbox',
  styleUrls: ['./textbox.component.css'],
  template:`
      <label for="filter">Filter: </label><input type="text" [formControl]="formControl">
      Value: {{ formControl.value }}
  `
})
export class TextboxComponent {
  @ContentChild(InteractionsDirective) interactionsDir?: InteractionsDirective;
  private interactions?: Interactions;

  formControl = new FormControl('');


  ngAfterContentInit(): void {
    // Ensure the user-specified interactions are supported.
    this.interactions = this.interactionsDir?.get();
    this.interactions?.checkForSupportedActions(supportedActions);

    this.formControl.valueChanges.subscribe(value => {
      this.interactions?.update(TEXTBOX, TEXT_CHANGED, new ValueMap(new Map([['value', new StringValue(value!)]])));
    });
  }

}
