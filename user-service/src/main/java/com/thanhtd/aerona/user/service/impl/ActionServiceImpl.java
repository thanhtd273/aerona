package com.thanhtd.aerona.user.service.impl;

import com.thanhtd.aerona.user.dao.ActionDao;
import com.thanhtd.aerona.user.dto.ActionInfo;
import com.thanhtd.aerona.user.model.Action;
import com.thanhtd.aerona.user.service.ActionService;
import com.thanhtd.aerona.base.constant.DataStatus;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;

import java.util.Date;
import java.util.List;

@Service
@RequiredArgsConstructor
public class ActionServiceImpl implements ActionService {

    private final ActionDao actionDao;

    @Override
    public List<Action> findAll() {
        return actionDao.findAll();
    }

    @Override
    public Action findByActionId(String actionId) throws LogicException {
        if (ObjectUtils.isEmpty(actionId))
            throw new LogicException(ErrorCode.ID_NULL);

        return actionDao.findByActionId(actionId);
    }

    @Override
    public Action findByCode(Integer code) throws LogicException {
        if (ObjectUtils.isEmpty(code))
            throw new LogicException(ErrorCode.NULL_VALUE);
        return actionDao.findByCode(code);
    }

    @Override
    public Action createAction(ActionInfo actionInfo) throws LogicException {
        if (ObjectUtils.isEmpty(actionInfo))
            throw new LogicException(ErrorCode.DATA_NULL);
        if (ObjectUtils.isEmpty(actionInfo.getName()))
            throw new LogicException(ErrorCode.BLANK_FIELD);

        Action action = new Action();
        action.setName(actionInfo.getName());
        if (!ObjectUtils.isEmpty(actionInfo.getDescription()))
            action.setDescription(actionInfo.getDescription());
        if (!ObjectUtils.isEmpty(actionInfo.getStatus()))
            action.setStatus(actionInfo.getStatus());

        action.setCreateDate(new Date());
        action.setStatus(DataStatus.ACTIVE);
        action = actionDao.save(action);

        return action;
    }

    @Override
    public Action updateAction(String actionId, ActionInfo actionInfo) throws LogicException {
        if (ObjectUtils.isEmpty(actionInfo) || ObjectUtils.isEmpty(actionId))
            throw new LogicException(ErrorCode.DATA_NULL);

        Action action = actionDao.findByActionId(actionId);
        if (ObjectUtils.isEmpty(action))
            throw new LogicException(ErrorCode.DATA_NULL);

        if (!ObjectUtils.isEmpty(actionInfo.getName()))
            action.setName(actionInfo.getName());
        if (!ObjectUtils.isEmpty(actionInfo.getDescription())) {
            action.setDescription(actionInfo.getDescription());
        }
        if (!ObjectUtils.isEmpty(actionInfo.getStatus())) {
            action.setStatus(actionInfo.getStatus());
        }
        action.setModifiedDate(new Date());
        action = actionDao.save(action);

        return action;
    }

    @Override
    public ErrorCode deleteAction(String actionId) {
        if (ObjectUtils.isEmpty(actionId))
            return ErrorCode.ID_NULL;
        Action action = actionDao.findByActionId(actionId);
        if (ObjectUtils.isEmpty(action))
            return ErrorCode.DATA_NULL;
        action.setModifiedDate(new Date());
        action.setStatus(DataStatus.DELETED);
        return ErrorCode.SUCCESS;
    }
}
