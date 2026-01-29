import { Injectable } from "@nestjs/common";
import { HttpService } from "@nestjs/axios";
import { firstValueFrom } from "rxjs";
import { AxiosRequestConfig } from "axios";

@Injectable()
export class ProxyService {
  constructor(private readonly httpService: HttpService) {}

  async proxyRequest(
    targetUrl: string,
    method: string,
    data?: any,
    headers?: any,
  ): Promise<any> {
    const config: AxiosRequestConfig = {
      method,
      url: targetUrl,
      headers: {
        ...headers,
        "Content-Type": "application/json",
      },
    };

    if (data && (method === "POST" || method === "PUT" || method === "PATCH")) {
      config.data = data;
    }

    try {
      const response = await firstValueFrom(this.httpService.request(config));
      return response.data;
    } catch (error) {
      if (error.response) {
        throw {
          statusCode: error.response.status,
          message: error.response.data.message || "Service error",
          error: error.response.data.error,
        };
      }
      throw error;
    }
  }
}
