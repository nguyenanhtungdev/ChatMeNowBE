import { Controller, All, Req, Res, HttpStatus } from "@nestjs/common";
import { Request, Response } from "express";
import { ProxyService } from "../services/proxy.service";

@Controller("api/auth")
export class AuthProxyController {
  private readonly authServiceUrl = process.env.AUTH_SERVICE_URL;

  constructor(private readonly proxyService: ProxyService) {}

  @All("*")
  async proxyAuth(@Req() req: Request, @Res() res: Response) {
    const path = req.url.replace("/api/auth", "");
    const targetUrl = `${this.authServiceUrl}/auth${path}`;

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
