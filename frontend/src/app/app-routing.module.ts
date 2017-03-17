import { AuthGuard } from './services/auth.guard';
import { LoginComponent } from './components/login/login.component';
import { HomeComponent } from './components/home/home.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LogoutComponent } from "./components/logout/logout.component";
import { SearchComponent } from "./components/search/search.component";
import { IntegrateComponent } from "./components/integrate/integrate.component";
import { CreateComponent } from "./components/integrate/create/create.component";

const appRoutes: Routes = [
    { path: 'home', component: HomeComponent, canActivate:[AuthGuard]},
    { path: 'login', component: LoginComponent },
    { path: 'logout', component: LogoutComponent },
    { path: 'search', component: SearchComponent },
    { path: 'integrate', component: IntegrateComponent },
    { path: 'integrate/create/:service', component: CreateComponent },
    { path: '', redirectTo: '/integrate', pathMatch: 'full' },
    { path: '**', component: PageNotFoundComponent }
];

@NgModule({
    imports: [
        RouterModule.forRoot(appRoutes)
    ],
    exports: [
        RouterModule
    ]
})
export class AppRoutingModule { }