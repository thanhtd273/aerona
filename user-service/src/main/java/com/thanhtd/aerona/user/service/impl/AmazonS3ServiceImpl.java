package com.thanhtd.aerona.user.service.impl;

import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.model.CannedAccessControlList;
import com.amazonaws.services.s3.model.DeleteObjectRequest;
import com.amazonaws.services.s3.model.PutObjectRequest;
import com.thanhtd.aerona.user.service.AmazonS3Service;
import com.thanhtd.aerona.base.exception.LogicException;
import jodd.util.StringPool;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.Date;
import java.util.Objects;

@Service
@RequiredArgsConstructor
public class AmazonS3ServiceImpl implements AmazonS3Service {
    @Value("${aws.s3.bucket.name}")
    private String s3BucketName;

    @Value("${aws.s3.endpoint.url}")
    private String s3EndpointUrl;

    private final AmazonS3 s3Client;

    @Override
    public String uploadFile(MultipartFile multipartFile) throws IOException, LogicException {
        File file = convertMultipartToFile(multipartFile);
        String fileName = generateFileName(multipartFile);
        String fileUrl = String.format("%s/%s", s3EndpointUrl, fileName);
        uploadFileToS3Bucket(fileName, file);
        return fileUrl;
    }

    @Override
    public String deleteFile(String fileUrl) {
        String fileName = fileUrl.substring(fileUrl.lastIndexOf(StringPool.SLASH) + 1);
        s3Client.deleteObject(new DeleteObjectRequest(s3BucketName, fileName));
        return String.format("Successfully delete file: %s", fileName) ;
    }

    private File convertMultipartToFile(MultipartFile multipartFile) throws IOException, LogicException {
        if (ObjectUtils.isEmpty(multipartFile))
            throw new LogicException("File is empty");
        File convFile = new File(Objects.requireNonNull(multipartFile.getOriginalFilename()));
        try (FileOutputStream fos = new FileOutputStream(convFile)) {
            fos.write(multipartFile.getBytes());
        }
        return convFile;
    }

    private String generateFileName(MultipartFile file) {
        String originFileName = file.getOriginalFilename();
        if (ObjectUtils.isEmpty(originFileName)) {
            return String.format("unnamed-file-%d", new Date().getTime()) ;
        }
        int dotIndex = originFileName.lastIndexOf(StringPool.DOT);
        String extension = originFileName.substring(dotIndex + 1);
        return String.format("%s-%d.%s",
                originFileName.substring(0, dotIndex).replace(StringPool.SPACE, StringPool.DASH), System.currentTimeMillis(), extension);
    }

    private void uploadFileToS3Bucket(String fileName, File file) {
        s3Client.putObject(new PutObjectRequest(s3BucketName, fileName, file).withCannedAcl(CannedAccessControlList.PublicRead));
    }
}
