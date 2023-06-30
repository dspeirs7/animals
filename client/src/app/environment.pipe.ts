import { Pipe, PipeTransform } from '@angular/core';
import { environment } from 'src/environments/environment';

@Pipe({
  name: 'env',
  standalone: true,
})
export class EnvironmentPipe implements PipeTransform {
  transform(value: keyof typeof environment): string | boolean {
    return environment[value];
  }
}
