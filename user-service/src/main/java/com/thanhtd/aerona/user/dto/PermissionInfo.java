package com.thanhtd.aerona.user.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class PermissionInfo implements Serializable {
    private String roleId;

    private String name;

    private String description;

    private String status;
}
