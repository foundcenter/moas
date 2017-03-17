import { AuthGuard } from './services/auth.guard';
import { MyAuthConfig } from './auth_config';
import { AppRoutingModule } from './app-routing.module';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { HomeComponent } from './components/home/home.component';
import { LoginComponent } from './components/login/login.component';
import { Ng2UiAuthModule } from "ng2-ui-auth/commonjs/ng2-ui-auth.module";
import { LogoutComponent } from './components/logout/logout.component';
import { SearchComponent } from "./components/search/search.component";
import { FocusDirective } from './focus.directive';
import { ResultComponent } from './components/search/result/result.component';
import { IntegrateComponent } from './components/integrate/integrate.component';
import { CreateComponent } from './components/integrate/create/create.component';

@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent,
    PageNotFoundComponent,
    HomeComponent,
    LoginComponent,
    LogoutComponent,
    SearchComponent,
    FocusDirective,
    ResultComponent,
    IntegrateComponent,
    CreateComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AppRoutingModule,
    Ng2UiAuthModule.forRoot(MyAuthConfig),
  ],
  providers: [
    AuthGuard
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
