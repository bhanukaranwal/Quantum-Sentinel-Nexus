// services/auth-gateway-svc/src/app.module.ts
import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';

@Module({
  imports: [
    // In a real application, you would import other modules here, for example:
    // - A ConfigModule to manage environment variables.
    // - An AuthModule to handle OIDC/JWT logic with Passport.js.
    // - A ThrottlerModule for rate limiting.
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
