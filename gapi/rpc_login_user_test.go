func TestLoginUser(t *testing.T) {
	user := randomUser()
	password := user.HashedPassword // assuming randomUser sets HashedPassword to the plain password for test

	testCases := []struct {
		name          string
		req           *pb.LoginUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, resp *pb.LoginUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(randomSession(user.Username), nil)
			},
			checkResponse: func(t *testing.T, resp *pb.LoginUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotEmpty(t, resp.AccessToken)
				require.NotEmpty(t, resp.RefreshToken)
				require.NotNil(t, resp.User)
				require.Equal(t, user.Username, resp.User.Username)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.LoginUserRequest{
				Username: "notfound",
				Password: "password",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq("notfound")).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, resp *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
				require.Nil(t, resp)
			},
		},
		{
			name: "IncorrectPassword",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: "wrongpassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, resp *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
				require.Nil(t, resp)
			},
		},
		{
			name: "InternalError",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, resp *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
				require.Nil(t, resp)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store, nil)
			resp, err := server.LoginUser(context.Background(), tc.req)
			tc.checkResponse(t, resp, err)
		})
	}
}

// randomUser returns a random user struct for testing.
func randomUser() db.User {
	password := util.RandomString(6)
	hashedPassword, _ := util.HashPassword(password)
	return db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
}

// randomSession returns a random session struct for testing.
func randomSession(username string) db.Session {
	return db.Session{
		ID:           util.RandomUUID(),
		Username:     username,
		RefreshToken: util.RandomString(32),
		UserAgent:    "test-agent",
		ClientIp:     "127.0.0.1",
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(time.Hour),
		CreatedAt:    time.Now(),
	}
}
