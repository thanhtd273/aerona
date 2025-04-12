package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.user.dto.OTPInfo;
import com.thanhtd.aerona.user.dto.UpdatePasswordInfo;
import com.thanhtd.aerona.user.dto.UserInfo;
import com.thanhtd.aerona.user.model.User;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;

public interface UserService {
    User createUser(UserInfo userInfo) throws LogicException;

    User findByUserId(String userId) throws LogicException;

    User findByEmail(String email) throws LogicException;

    UserInfo getUserInfo(User user) throws LogicException;

    UserInfo getUserInfoByToken(String token) throws LogicException;

    UserInfo getCurrentUser() throws LogicException;

    User updateUser(String userId, UserInfo userInfo) throws LogicException;

    ErrorCode activateUser(OTPInfo activationInfo) throws LogicException;

    ErrorCode deactivateUser(String userId) throws LogicException;

    ErrorCode changePassword(UpdatePasswordInfo updatePasswordInfo) throws LogicException;

    ErrorCode forgotPassword(UserInfo userInfo) throws LogicException;

    ErrorCode verifyPasswordResetCode(OTPInfo passwordResetInfo) throws LogicException;

    ErrorCode addRole(UserInfo userInfo) throws LogicException;

    ErrorCode removeRole(UserInfo userInfo) throws LogicException;

}
