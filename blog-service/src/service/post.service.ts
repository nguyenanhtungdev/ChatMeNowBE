import {
  Injectable,
  NotFoundException,
  ForbiddenException,
} from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { BlogPost } from "../entities/post.entity";
import { CreatePostDto, UpdatePostDto } from "../dto/post.dto";

@Injectable()
export class PostService {
  constructor(
    @InjectRepository(BlogPost)
    private postRepo: Repository<BlogPost>,
  ) {}

  async create(userId: string, dto: CreatePostDto) {
    const post = this.postRepo.create({
      ...dto,
      userId,
    });
    return this.postRepo.save(post);
  }

  async findAll(userId?: string, status?: string) {
    const query = this.postRepo.createQueryBuilder("post");

    if (userId) {
      query.where("post.user_id = :userId", { userId });
    }

    if (status) {
      query.andWhere("post.status = :status", { status });
    } else {
      // Default: only show published posts
      query.andWhere("post.status = :status", { status: "published" });
    }

    query.orderBy("post.created_at", "DESC");

    return query.getMany();
  }

  async findOne(id: string) {
    const post = await this.postRepo.findOne({ where: { id } });
    if (!post) {
      throw new NotFoundException("Post not found");
    }

    // Increment view count
    post.viewCount += 1;
    await this.postRepo.save(post);

    return post;
  }

  async update(id: string, userId: string, dto: UpdatePostDto) {
    const post = await this.postRepo.findOne({ where: { id } });

    if (!post) {
      throw new NotFoundException("Post not found");
    }

    if (post.userId !== userId) {
      throw new ForbiddenException("You can only update your own posts");
    }

    Object.assign(post, dto);
    return this.postRepo.save(post);
  }

  async publish(id: string, userId: string) {
    const post = await this.postRepo.findOne({ where: { id } });

    if (!post) {
      throw new NotFoundException("Post not found");
    }

    if (post.userId !== userId) {
      throw new ForbiddenException("You can only publish your own posts");
    }

    post.status = "published";
    post.publishedAt = new Date();
    return this.postRepo.save(post);
  }

  async delete(id: string, userId: string) {
    const post = await this.postRepo.findOne({ where: { id } });

    if (!post) {
      throw new NotFoundException("Post not found");
    }

    if (post.userId !== userId) {
      throw new ForbiddenException("You can only delete your own posts");
    }

    await this.postRepo.remove(post);
    return { message: "Post deleted successfully" };
  }
}
