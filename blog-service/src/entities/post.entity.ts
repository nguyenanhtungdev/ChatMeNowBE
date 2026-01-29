import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
} from "typeorm";

@Entity("blog_posts")
export class BlogPost {
  @PrimaryGeneratedColumn("uuid")
  id: string;

  @Column({ name: "user_id" })
  userId: string;

  @Column({ length: 500 })
  title: string;

  @Column("text")
  content: string;

  @Column({ nullable: true, type: "text" })
  excerpt?: string;

  @Column({ name: "cover_image", nullable: true })
  coverImage?: string;

  @Column({ default: "draft" })
  status: "draft" | "published" | "archived";

  @Column("text", { array: true, nullable: true })
  tags?: string[];

  @Column({ name: "view_count", default: 0 })
  viewCount: number;

  @Column({ name: "published_at", nullable: true })
  publishedAt?: Date;

  @CreateDateColumn({ name: "created_at" })
  createdAt: Date;

  @UpdateDateColumn({ name: "updated_at" })
  updatedAt: Date;
}
