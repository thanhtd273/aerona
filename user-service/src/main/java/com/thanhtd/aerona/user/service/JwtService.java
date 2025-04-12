package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.user.model.User;
import com.thanhtd.aerona.base.constant.ErrorCode;
import io.jsonwebtoken.Claims;

import java.util.Map;
import java.util.function.Function;

public interface JwtService {
    String extractEmail(String token);

    <T> T extractClaim(String token, Function<Claims, T> claimsSolver);

    String generateToken(User user);

    String generateToken(Map<String, Object> extraClaims, User user);

    ErrorCode validateToken(String token, User user);

    int getExpireIn();
}
