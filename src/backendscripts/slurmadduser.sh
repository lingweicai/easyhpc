# Step 1: make sure account exists
sacctmgr add account dev Description="Developer Account"

# Step 2: add user to account (association)
sacctmgr add user dev account=dev