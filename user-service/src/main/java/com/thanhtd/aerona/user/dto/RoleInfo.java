package com.thanhtd.aerona.user.dto;

import lombok.Data;

@Data
public class RoleInfo {
    private String roleId;

    private String name;

    private String description;

    private String status;

    private String permissions;
}
