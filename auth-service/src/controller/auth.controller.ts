import {
  Controller,
  Post,
  Get,
  Body,
  Headers,
  Req,
  UseGuards,
} from "@nestjs/common";
import { AuthService } from "../service/auth.service";
import { RegisterDto, LoginDto, RefreshTokenDto } from "../dto/auth.dto";
import { JwtAuthGuard } from "../guards/jwt-auth.guard";

@Controller("auth")
export class AuthController {
  constructor(private authService: AuthService) {}

  @Post("register")
  register(@Body() dto: RegisterDto) {
    return this.authService.register(dto);
  }

  @Post("login")
  login(
    @Body() dto: LoginDto,
    @Headers("x-forwarded-for") ip: string,
    @Headers("user-agent") userAgent: string,
  ) {
    return this.authService.login(dto, ip, userAgent);
  }

  @Post("refresh")
  refresh(@Body() dto: RefreshTokenDto) {
    return this.authService.refresh(dto.refreshToken);
  }

  @Get("me")
  @UseGuards(JwtAuthGuard)
  getMe(@Req() req: any) {
    return this.authService.getMe(req.user.sub);
  }

  @Post("logout")
  @UseGuards(JwtAuthGuard)
  logout(@Req() req: any, @Body() dto: RefreshTokenDto) {
    return this.authService.logout(req.user.sub, dto.refreshToken);
  }

  @Get("health")
  health() {
    return { status: "ok", service: "auth-service" };
  }
}
