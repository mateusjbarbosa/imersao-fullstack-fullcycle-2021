import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ModelNotFoundException } from './exceptions/model-not.found.exception';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  app.setGlobalPrefix('api/v1');

  app.useGlobalFilters(new ModelNotFoundException());

  await app.listen(3000);
}
bootstrap();
