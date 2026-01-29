import { Module } from "@nestjs/common";
import { TypeOrmModule } from "@nestjs/typeorm";
import { JwtModule } from "@nestjs/jwt";
import { PostController } from "./controller/post.controller";
import { PostService } from "./service/post.service";
import { BlogPost } from "./entities/post.entity";
import { JwtAuthGuard } from "./guards/jwt-auth.guard";

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: "postgres",
      url: process.env.DATABASE_URL,
      entities: [BlogPost],
      synchronize: false,
      logging: process.env.NODE_ENV !== "production",
    }),
    TypeOrmModule.forFeature([BlogPost]),
    JwtModule.register({
      secret: process.env.JWT_SECRET || "your-secret-key",
    }),
  ],
  controllers: [PostController],
  providers: [PostService, JwtAuthGuard],
})
export class AppModule {}
