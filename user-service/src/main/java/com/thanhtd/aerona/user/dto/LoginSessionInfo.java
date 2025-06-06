package com.thanhtd.aerona.user.dto;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class LoginSessionInfo {
    private String sessionId;

    private String token;

    private UserInfo userInfo;

    private Integer expireIn;
}
