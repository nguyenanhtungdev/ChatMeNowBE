import {
  Injectable,
  UnauthorizedException,
  ConflictException,
} from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { JwtService } from "@nestjs/jwt";
import * as bcrypt from "bcrypt";
import { User } from "../entities/user.entity";
import { RefreshToken } from "../entities/refresh-token.entity";
import { RegisterDto, LoginDto } from "../dto/auth.dto";

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User)
    private userRepo: Repository<User>,
    @InjectRepository(RefreshToken)
    private refreshTokenRepo: Repository<RefreshToken>,
    private jwtService: JwtService,
  ) {}

  async register(dto: RegisterDto) {
    // Check if user exists
    const existing = await this.userRepo.findOne({
      where: [{ email: dto.email }, { username: dto.username }],
    });

    if (existing) {
      throw new ConflictException("Username or email already exists");
    }

    // Hash password
    const passwordHash = await bcrypt.hash(dto.password, 10);

    // Create user
    const user = this.userRepo.create({
      username: dto.username,
      email: dto.email,
      passwordHash,
    });

    await this.userRepo.save(user);

    // Return user info without tokens (user must login to get tokens)
    return {
      message: "Registration successful. Please login to continue.",
      user: {
        id: user.id,
        username: user.username,
        email: user.email,
      },
    };
  }

  async login(dto: LoginDto, ipAddress?: string, userAgent?: string) {
    // Find user
    const user = await this.userRepo.findOne({
      where: [
        { email: dto.emailOrUsername },
        { username: dto.emailOrUsername },
      ],
    });

    if (!user || !user.isActive) {
      throw new UnauthorizedException("Invalid credentials");
    }

    // Verify password
    const valid = await bcrypt.compare(dto.password, user.passwordHash);
    if (!valid) {
      throw new UnauthorizedException("Invalid credentials");
    }

    // Generate tokens
    const tokens = await this.generateTokens(user);

    // Save refresh token
    const expiresAt = new Date();
    expiresAt.setDate(expiresAt.getDate() + 7); // 7 days

    const tokenHash = await bcrypt.hash(tokens.refreshToken, 10);

    await this.refreshTokenRepo.save({
      userId: user.id,
      tokenHash,
      deviceId: dto.deviceId,
      deviceName: dto.deviceName,
      ipAddress,
      userAgent,
      expiresAt,
    });

    return tokens;
  }

  async refresh(refreshToken: string) {
    try {
      const payload = this.jwtService.verify(refreshToken);

      // Find refresh token in DB
      const tokens = await this.refreshTokenRepo.find({
        where: { userId: payload.sub },
      });

      let valid = false;
      for (const token of tokens) {
        if (await bcrypt.compare(refreshToken, token.tokenHash)) {
          if (token.expiresAt < new Date()) {
            throw new UnauthorizedException("Refresh token expired");
          }
          valid = true;
          break;
        }
      }

      if (!valid) {
        throw new UnauthorizedException("Invalid refresh token");
      }

      const user = await this.userRepo.findOne({ where: { id: payload.sub } });
      if (!user) {
        throw new UnauthorizedException("User not found");
      }

      return this.generateTokens(user);
    } catch (error) {
      throw new UnauthorizedException("Invalid refresh token");
    }
  }

  async getMe(userId: string) {
    const user = await this.userRepo.findOne({ where: { id: userId } });
    if (!user) {
      throw new UnauthorizedException("User not found");
    }

    const { passwordHash, ...result } = user;
    return result;
  }

  async logout(userId: string, refreshToken: string) {
    const tokens = await this.refreshTokenRepo.find({ where: { userId } });

    for (const token of tokens) {
      if (await bcrypt.compare(refreshToken, token.tokenHash)) {
        await this.refreshTokenRepo.remove(token);
        break;
      }
    }

    return { message: "Logged out successfully" };
  }

  private async generateTokens(user: User) {
    const payload = {
      sub: user.id,
      username: user.username,
      email: user.email,
    };

    const accessToken = this.jwtService.sign(payload);
    const refreshToken = this.jwtService.sign(payload, { expiresIn: "7d" });

    return {
      accessToken,
      refreshToken,
      user: {
        id: user.id,
        username: user.username,
        email: user.email,
        avatarUrl: user.avatarUrl,
      },
    };
  }
}
