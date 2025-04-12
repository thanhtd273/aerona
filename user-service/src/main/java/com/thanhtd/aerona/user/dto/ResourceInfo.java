package com.thanhtd.aerona.user.dto;

import lombok.Data;

@Data
public class ResourceInfo {
    private String resourceId;

    private String name;

    private String path;

    private String description;

    private String status;
}
