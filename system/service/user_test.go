package service

// func TestUser(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
//
// 	usr := &types.User{ID: factory.Sonyflake.NextID()}
//
// 	usrRpoMock := NewMockRepository(mockCtrl)
// 	usrRpoMock.EXPECT().WithCtx(gomock.Any()).AnyTimes().Return(usrRpoMock)
// 	usrRpoMock.EXPECT().
// 		FindUserByID(usr.ID).
// 		Times(1).
// 		Return(usr, nil)
//
// 	svc := User()
// 	svc.rpo = usrRpoMock
//
// 	found, err := svc.FindByID(context.Background(), usr.ID)
// 	if err != nil {
// 		t.Fatal("Did not expect an error")
// 	}
//
// 	if found == nil {
// 		t.Fatal("Expecting an user to be found")
// 	}
//
// 	if found.ID != usr.ID {
// 		t.Fatal("Expecting found user to have the same ID as the find param")
// 	}
// }
