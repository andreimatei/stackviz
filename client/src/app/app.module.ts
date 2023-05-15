import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatCardModule } from '@angular/material/card';
import { MatPaginatorModule } from '@angular/material/paginator';
import { CoreModule, DataTableModule, WeightedTreeModule, TextFieldModule } from 'traceviz/dist/ngx-traceviz-lib';
import { AppComponent } from './app.component';
import { ComponentsModule } from "./components/components.module";

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    CoreModule,
    DataTableModule,
    HttpClientModule,
    MatCardModule,
    MatPaginatorModule,
    WeightedTreeModule,
    ComponentsModule,
    TextFieldModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
