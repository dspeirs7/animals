import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideAnimations } from '@angular/platform-browser/animations';
import { ChickenService } from './chicken.service';
import { HttpClientModule } from '@angular/common/http';
import { EnvironmentPipe } from './environment.pipe';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideAnimations(),
    { provide: ChickenService, useClass: ChickenService },
    { provide: EnvironmentPipe, useClass: EnvironmentPipe },
    importProvidersFrom(HttpClientModule),
  ],
};
