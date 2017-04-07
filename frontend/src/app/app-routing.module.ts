import { AuthGuard } from './services/auth.guard';
import { LoginComponent } from './components/login/login.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LogoutComponent } from './components/logout/logout.component';
import { SearchComponent } from './components/search/search.component';
import { ConnectComponent } from './components/connect/connect.component';

const appRoutes: Routes = [
    { path: 'login', component: LoginComponent },
    { path: 'logout', component: LogoutComponent, canActivate:[AuthGuard]},
    { path: 'search', component: SearchComponent, canActivate:[AuthGuard]},
    { path: 'connect', component: ConnectComponent, canActivate:[AuthGuard]},
    { path: '', redirectTo: '/search', pathMatch: 'full' },
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