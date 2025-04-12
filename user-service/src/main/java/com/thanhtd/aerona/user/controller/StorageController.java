package com.thanhtd.aerona.user.controller;

import com.thanhtd.aerona.user.service.AmazonS3Service;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.base.exception.ExceptionHandler;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RequestPart;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/api/v1/user-service/storage")
@RequiredArgsConstructor
public class StorageController {

    private static final Logger logger = LoggerFactory.getLogger(StorageController.class);

    private final AmazonS3Service amazonS3Service;

    @PostMapping("/upload")
    public APIResponse upload(HttpServletResponse response, @RequestPart(value = "file") MultipartFile file) {
        long start = System.currentTimeMillis();
        try {
            String fileUrl = amazonS3Service.uploadFile(file);
            logger.info("Call API POST /api/v1/user-service/storage/upload success, file {}", file.getOriginalFilename());
            return new APIResponse(ErrorCode.SUCCESS, "Upload file success", System.currentTimeMillis() - start, fileUrl);
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/storage/upload failed, file {}, error: {}", file.getOriginalFilename(), e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @DeleteMapping("/delete")
    public APIResponse deleteFile(HttpServletResponse response, @RequestParam("fileUrl") String fileUrl) {
        long start = System.currentTimeMillis();
        try {
            String message = amazonS3Service.deleteFile(fileUrl);
            logger.info("Call API DELETE /api/v1/user-service/storage/delete?fileUr={} success", fileUrl);
            return new APIResponse(ErrorCode.SUCCESS, "Delete file success", System.currentTimeMillis() - start, message);
        } catch (Exception e) {
            logger.error("Call API DELETE /api/v1/user-service/storage/delete?fileUrl={} failed, error: {}", fileUrl, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }
}
