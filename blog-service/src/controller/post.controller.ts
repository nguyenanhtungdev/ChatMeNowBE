import {
  Controller,
  Get,
  Post,
  Put,
  Delete,
  Patch,
  Body,
  Param,
  Query,
  Req,
  UseGuards,
} from "@nestjs/common";
import { PostService } from "../service/post.service";
import { CreatePostDto, UpdatePostDto } from "../dto/post.dto";
import { JwtAuthGuard } from "../guards/jwt-auth.guard";

@Controller("posts")
export class PostController {
  constructor(private postService: PostService) {}

  @Get("health")
  health() {
    return { status: "ok", service: "blog-service" };
  }

  @Post()
  @UseGuards(JwtAuthGuard)
  create(@Req() req: any, @Body() dto: CreatePostDto) {
    return this.postService.create(req.user.sub, dto);
  }

  @Get()
  findAll(@Query("userId") userId?: string, @Query("status") status?: string) {
    return this.postService.findAll(userId, status);
  }

  @Get(":id")
  findOne(@Param("id") id: string) {
    return this.postService.findOne(id);
  }

  @Put(":id")
  @UseGuards(JwtAuthGuard)
  update(@Param("id") id: string, @Req() req: any, @Body() dto: UpdatePostDto) {
    return this.postService.update(id, req.user.sub, dto);
  }

  @Patch(":id/publish")
  @UseGuards(JwtAuthGuard)
  publish(@Param("id") id: string, @Req() req: any) {
    return this.postService.publish(id, req.user.sub);
  }

  @Delete(":id")
  @UseGuards(JwtAuthGuard)
  delete(@Param("id") id: string, @Req() req: any) {
    return this.postService.delete(id, req.user.sub);
  }
}
