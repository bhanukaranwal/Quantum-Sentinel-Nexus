// services/auth-gateway-svc/src/app.controller.ts
import { Controller, Get } from '@nestjs/common';
import { AppService } from './app.service';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  /**
   * Defines a GET endpoint at the root of the service.
   * This is useful as a health check to confirm the gateway is running.
   */
  @Get()
  getHealthStatus(): { status: string } {
    return this.appService.getHealthStatus();
  }

  // In a real gateway, you might have other routes here, for example:
  // - A '/login' route to initiate the OIDC flow.
  // - A '/userinfo' route that is protected and returns user details after authentication.
}
