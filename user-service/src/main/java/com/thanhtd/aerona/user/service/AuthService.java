package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.user.dto.Credential;
import com.thanhtd.aerona.user.dto.LoginSessionInfo;
import com.thanhtd.aerona.user.dto.UserInfo;
import com.thanhtd.aerona.user.model.User;
import com.thanhtd.aerona.base.exception.LogicException;
import jakarta.servlet.http.HttpServletRequest;

public interface AuthService {
    UserInfo verifyToken(String token) throws LogicException;

    LoginSessionInfo login(HttpServletRequest request, Credential credential) throws LogicException;

    User signUp(UserInfo userInfo) throws LogicException;
}
