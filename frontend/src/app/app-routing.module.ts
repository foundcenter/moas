import { AuthGuard } from './services/auth.guard';
import { LoginComponent } from './components/login/login.component';
import { HomeComponent } from './components/home/home.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LogoutComponent } from "./components/logout/logout.component";
import { SearchComponent } from "./components/search/search.component";
import { IntegrateComponent } from "./components/integrate/integrate.component";

const appRoutes: Routes = [
    { path: 'login', component: LoginComponent },
    { path: 'home', component: HomeComponent, canActivate:[AuthGuard]},
    { path: 'logout', component: LogoutComponent, canActivate:[AuthGuard]},
    { path: 'search', component: SearchComponent, canActivate:[AuthGuard]},
    { path: 'integrate', component: IntegrateComponent, canActivate:[AuthGuard]},
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