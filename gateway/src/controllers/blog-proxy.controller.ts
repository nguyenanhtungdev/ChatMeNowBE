import {
  Controller,
  All,
  Req,
  Res,
  HttpStatus,
  UseGuards,
} from "@nestjs/common";
import { Request, Response } from "express";
import { ProxyService } from "../services/proxy.service";
import { JwtAuthGuard } from "../guards/jwt-auth.guard";

@Controller("api/blog")
@UseGuards(JwtAuthGuard)
export class BlogProxyController {
  private readonly blogServiceUrl = process.env.BLOG_SERVICE_URL;

  constructor(private readonly proxyService: ProxyService) {}

  @All("*")
  async proxyBlog(@Req() req: Request, @Res() res: Response) {
    const path = req.url.replace("/api/blog", "");
    const targetUrl = `${this.blogServiceUrl}/posts${path}`;

    try {
      const result = await this.proxyService.proxyRequest(
        targetUrl,
        req.method,
        req.body,
        req.headers,
      );
      return res.status(HttpStatus.OK).json(result);
    } catch (error) {
      return res
        .status(error.statusCode || HttpStatus.INTERNAL_SERVER_ERROR)
        .json(error);
    }
  }
}
