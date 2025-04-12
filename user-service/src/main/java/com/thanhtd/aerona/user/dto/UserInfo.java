package com.thanhtd.aerona.user.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class UserInfo implements Serializable {
    private String userId;

    private String name;

    private String email;

    private String password;

    private String confirmPassword;

    private String avatar;

    private String roles;

    private String status;
}
