package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.user.dto.ActionInfo;
import com.thanhtd.aerona.user.model.Action;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;

import java.util.List;

public interface ActionService {
    List<Action> findAll();

    Action findByActionId(String actionId) throws LogicException;

    Action findByCode(Integer code) throws LogicException;

    Action createAction(ActionInfo actionInfo) throws LogicException;

    Action updateAction(String actionId, ActionInfo actionInfo) throws LogicException;

    ErrorCode deleteAction(String actionId);
}
