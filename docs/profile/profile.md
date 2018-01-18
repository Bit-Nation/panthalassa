# Profile

`profile/profile.js` is used to interact with the profiles saved on the device.
Profiles are saved in a realm database.
You will might see methods that relate to `PublicProfile`.
The public profile is the "normal" profile with some additonal fields.
Such as e.g. your ethereum addresses. The public profile is not safed.
It's created on runtime based on your profile + some other data.