package com.thanhtd.aerona.user.service.impl;

import com.thanhtd.aerona.user.dto.Credential;
import com.thanhtd.aerona.user.dto.LoginSessionInfo;
import com.thanhtd.aerona.user.dto.UserInfo;
import com.thanhtd.aerona.user.model.User;
import com.thanhtd.aerona.user.service.AuthService;
import com.thanhtd.aerona.user.service.JwtService;
import com.thanhtd.aerona.user.service.UserService;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Lazy;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;

@Service
@RequiredArgsConstructor
public class AuthServiceImpl implements AuthService {

    private final JwtService jwtService;

    private final UserService userService;

    @Lazy
    private final AuthenticationManager authenticationManager;

    @Override
    public UserInfo verifyToken(String token) throws LogicException {
        if (ObjectUtils.isEmpty(token) || !token.toLowerCase().startsWith("bearer")) {
            throw  new LogicException(ErrorCode.INVALID_TOKEN);
        }
        token = token.substring(7);
        String email = jwtService.extractEmail(token);
        User user = userService.findByEmail(email);
        ErrorCode tokenValidating = jwtService.validateToken(token, user);
        if (tokenValidating == ErrorCode.INVALID_TOKEN)
            throw new LogicException(ErrorCode.INVALID_TOKEN);
        return userService.getUserInfo(user);
    }

    @Override
    public LoginSessionInfo login(HttpServletRequest request, Credential credential) throws LogicException {
        User loginUser = userService.findByEmail(credential.getEmail());
        if (ObjectUtils.isEmpty(loginUser)) {
            throw new LogicException(ErrorCode.NOT_FOUND_USER);
        }
        Authentication usernamePasswordAuthentication = new UsernamePasswordAuthenticationToken(credential.getEmail(), credential.getPassword());
        Authentication authentication = authenticationManager.authenticate(usernamePasswordAuthentication);
        User user = (User) authentication.getPrincipal();
        if (ObjectUtils.isEmpty(user))
            throw new LogicException(ErrorCode.LOGIN_FAIL);

        UserInfo userInfo = userService.getUserInfo(user);
        request.getSession().setAttribute("user", userInfo);

        LoginSessionInfo loginSessionInfo = new LoginSessionInfo();
        loginSessionInfo.setSessionId(request.getSession().getId());
        loginSessionInfo.setUserInfo(userInfo);
        loginSessionInfo.setToken(jwtService.generateToken(user));
        loginSessionInfo.setExpireIn(jwtService.getExpireIn());

        return loginSessionInfo;
    }

    @Override
    public User signUp(UserInfo userInfo) throws LogicException {
        return userService.createUser(userInfo);
    }
}
