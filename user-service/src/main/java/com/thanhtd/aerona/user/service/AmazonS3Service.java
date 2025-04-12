package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.base.exception.LogicException;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;

public interface AmazonS3Service {
    String uploadFile(MultipartFile file) throws IOException, LogicException;

    String deleteFile(String fileUrl);
}
