import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatCardModule } from '@angular/material/card';
import { MatPaginatorModule } from '@angular/material/paginator';
import { CoreModule, DataTableModule, WeightedTreeModule, TextFieldModule } from 'traceviz/dist/ngx-traceviz-lib';
import { AppComponent } from './app.component';
import { ComponentsModule } from "./components/components.module";
import { StacksComponent } from "./components/stacks/stacks.component";
import { GraphQLModule } from './graphql/graphql.module';

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
    StacksComponent,
    GraphQLModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
