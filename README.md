# PhotoSpace

## API
### Build app
`make build`

### Start App
`make run`

### Run tests
`make test`


### Development
#### Hot reload
`air`
##### Installation
`https://github.com/cosmtrek/air`
`curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh -b "./airt" | sh -s -- -b ./utils`

##### Troubleshooting
Delete bound address
`lsof -ni tcp:8000`
`kill {pid}`