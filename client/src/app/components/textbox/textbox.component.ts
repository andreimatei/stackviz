import {Component, ContentChild} from '@angular/core';

import { InteractionsDirective } from 'traceviz/dist/ngx-traceviz-lib';
import { Interactions } from 'traceviz-core';

const TEXTBOX = 'textbox';
const TEXT_CHANGED = 'text_changed';

const supportedActions = new Array<[string, string]>(
    [TEXTBOX, TEXT_CHANGED],
);

@Component({
  selector: 'textbox',
  templateUrl: './textbox.component.html',
  styleUrls: ['./textbox.component.css']
})
export class TextboxComponent {
  @ContentChild(InteractionsDirective) interactionsDir?: InteractionsDirective;
  private interactions?: Interactions;
  
  ngAfterContentInit(): void {
    // Ensure the user-specified interactions are supported.
    this.interactions = this.interactionsDir?.get();
    this.interactions?.checkForSupportedActions(supportedActions);
  }

}
