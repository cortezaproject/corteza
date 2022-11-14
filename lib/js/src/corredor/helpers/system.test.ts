import { describe, it, beforeEach } from 'mocha'
import { expect } from 'chai'
import { stubObject, StubbedInstance } from 'ts-sinon'
import SystemHelper from './system'
import { User, Role } from '../../system/'
import { System as SystemAPI } from '../../api-clients'
import { kv } from './shared'

describe(__filename, () => {
  describe('supporting functions', () => {
    describe('resolving user', () => {
      let h: StubbedInstance<SystemHelper>

      beforeEach(() => {
        h = stubObject<SystemHelper>(
          new SystemHelper({ SystemAPI: new SystemAPI({}) }),
          [
            'findUsers',
            'findUserByID',
            'findUserByEmail',
            'findUserByHandle',
          ],
        )
      })

      it('should return first valid user without invoking the API', async () => {
        const u = new User({ userID: '444' })
        expect(await h.resolveUser(undefined, null, false, 0, '', u)).to.deep.equal(u)
        expect(h.findUserByID.notCalled, 'findUserByID call not expected').true
        expect(h.findUserByEmail.notCalled, 'findUserByEmail call not expected').true
        expect(h.findUserByHandle.notCalled, 'findUserByHandle call not expected').true
      })

      it('should resolve user ID and invoke the API', async () => {
        const u = new User({ userID: '444' })
        h.findUserByID.resolves(u)
        expect(await h.resolveUser(u.userID)).to.deep.equal(u)
        expect(h.findUserByID.calledOnceWith(u.userID), 'findUserByID call expected').true
        expect(h.findUserByEmail.notCalled, 'findUserByEmail call not expected').true
        expect(h.findUserByHandle.notCalled, 'findUserByHandle call not expected').true
      })

      it('should resolve handle', async () => {
        const u = new User({ handle: 'user-handle' })

        h.findUserByHandle.resolves(u)
        expect(await h.resolveUser(u.handle)).to.deep.equal(u)
        expect(h.findUserByHandle.calledOnceWith(u.handle), 'findUserByHandle call expected').true
        expect(h.findUserByID.notCalled, 'findUserByID call not expected').true
        expect(h.findUserByEmail.notCalled, 'findUserByEmail call not expected').true
      })

      it('should resolve numeric handle (first by ID then fallback to handle)', async () => {
        const u = new User({ handle: '42' })

        h.findUserByID.rejects()
        h.findUserByHandle.resolves(u)
        expect(await h.resolveUser(u.handle)).to.deep.equal(u)
        expect(h.findUserByID.calledOnceWith(u.handle), 'findUserByID call expected').true
        expect(h.findUserByHandle.calledOnceWith(u.handle), 'findUserByHandle call expected').true
        expect(h.findUserByEmail.notCalled, 'findUserByEmail call not expected').true
      })

      it('should resolve numeric handle (first by ID then fallback to handle)', async () => {
        const u = new User({ email: 'foo@bar.baz' })

        h.findUserByEmail.resolves(u)
        expect(await h.resolveUser(u.email)).to.deep.equal(u)
        expect(h.findUserByEmail.calledOnceWith(u.email), 'findUserByEmail call expected').true
        expect(h.findUserByID.notCalled, 'findUserByID call not expected').true
        expect(h.findUserByHandle.notCalled, 'findUserByHandle call not expected').true
      })
    })

    it.skip('resolving role')
  })

  describe('helpers', () => {
    let h: SystemHelper
    let systemApiStub: StubbedInstance<SystemAPI>

    beforeEach(() => {
      systemApiStub = stubObject<SystemAPI>(new SystemAPI({}))

      h = new SystemHelper({
        SystemAPI: systemApiStub,
      })
    })

    describe('user finding', () => {
      it('handles string filter', async () => {
        systemApiStub.userList.resolves({ set: [new User()] })
        await h.findUsers('filter')
        expect(systemApiStub.userList.calledOnceWith({ query: 'filter' }))
      })

      it('returns valid object', async () => {
        systemApiStub.userList.resolves({ set: [new User()] })
        expect((await h.findUsers()).set[0]).is.instanceOf(User)
      })
    })

    describe('user finding by ID', () => {
      it('returns valid object', async () => {
        systemApiStub.userRead.resolves(kv(new User()))
        expect(await h.findUserByID('1234')).is.instanceOf(User)
      })
    })

    describe('user saving', () => {
      it('should create new', async () => {
        const u = new User()

        systemApiStub.userCreate.resolves(kv(u))

        await h.saveUser(u)

        expect(systemApiStub.userCreate.calledOnceWith(kv(u)))
      })

      it('should update existing', async () => {
        const u = new User({ userID: '555' })

        systemApiStub.userUpdate.resolves(kv(u))

        await h.saveUser(u)

        expect(systemApiStub.userUpdate.calledOnceWith(kv(u))).true
      })
    })

    describe('user deleting', () => {
      it('should delete existing user', async () => {
        const u = new User({ userID: '222' })

        systemApiStub.userDelete.resolves(kv(u))

        await h.deleteUser(u)

        expect(systemApiStub.userDelete.calledOnceWith({ userID: u.userID })).true
      })

      it('should not delete fresh user', async () => {
        const user = new User()

        systemApiStub.userDelete.resolves()

        await h.deleteUser(user)

        expect(systemApiStub.userDelete.notCalled).true
      })
    })

    describe('roles finding', () => {
      it('handles string filter', async () => {
        systemApiStub.roleList.resolves({ set: [new Role()] })
        await h.findRoles('filter')
        expect(systemApiStub.roleList.calledOnceWith({ query: 'filter' }))
      })

      it('returns valid object', async () => {
        systemApiStub.roleList.resolves({ set: [new Role()] })
        expect((await h.findRoles()).set[0]).is.instanceOf(Role)
      })
    })

    describe('roles finding by id', () => {
      it('returns valid object', async () => {
        systemApiStub.roleRead.resolves(kv(new Role()))

        expect(await h.findRoleByID('1234')).is.instanceOf(Role)
        expect(systemApiStub.roleRead.calledOnce).true
      })
    })

    describe('role saving', () => {
      it('should create new', async () => {
        const role = new Role()

        systemApiStub.roleCreate.resolves(kv(role))

        await h.saveRole(role)

        expect(systemApiStub.roleCreate.calledOnceWith(kv(role))).true
      })

      it('should update existing', async () => {
        const role = new Role({ roleID: '555' })

        systemApiStub.roleUpdate.resolves(kv(role))

        await h.saveRole(role)

        expect(systemApiStub.roleUpdate.calledOnceWith(kv(role))).true
      })
    })

    describe('role deleting', () => {
      it('should delete existing role', async () => {
        const role = new Role({ roleID: '222' })

        systemApiStub.roleDelete.resolves()

        await h.deleteRole(role)

        expect(systemApiStub.roleDelete.calledOnceWith({ roleID: '222' })).true
      })

      it('should not delete fresh role', async () => {
        const role = new Role()

        systemApiStub.roleDelete.resolves()

        await h.deleteRole(role)

        expect(systemApiStub.roleDelete.notCalled).true
      })
    })

    describe('role membership ', () => {
      it('should add user to given role', async () => {
        const role = new Role({ roleID: '222' })
        const user = new User({ userID: '444' })

        systemApiStub.roleMemberAdd.resolves()

        await h.addUserToRole(user, role)
        expect(systemApiStub.roleMemberAdd.calledOnceWith({ roleID: '222', userID: '444' })).true
      })

      it('should remove user from given role', async () => {
        const role = new Role({ roleID: '222' })
        const user = new User({ userID: '444' })

        systemApiStub.roleMemberRemove.resolves()

        await h.removeUserFromRole(user, role)

        expect(systemApiStub.roleMemberRemove.calledOnceWith({ roleID: '222', userID: '444' })).true
      })
    })
  })
})
