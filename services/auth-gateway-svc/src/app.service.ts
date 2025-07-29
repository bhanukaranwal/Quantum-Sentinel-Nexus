// services/auth-gateway-svc/src/app.service.ts
import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  /**
   * Returns a simple status object to indicate the service is running.
   * @returns An object with a 'status' key.
   */
  getHealthStatus(): { status: string } {
    return { status: 'Auth Gateway Service is running' };
  }
}
