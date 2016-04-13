package main

type apiTestCase struct {
	requestBody string
	pipeline    Pipeline
}

var apiTestCases = []apiTestCase{
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "Step Name",
	      "image": "ubuntu:14.04",
	      "cmds": [
	        "ls -la",
	        "touch hello.txt",
	        "ls -la"
	      ]
			}
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "Step Name",
					ImageName: "ubuntu:14.04",
					Cmds: []Cmd{
						"ls -la",
						"touch hello.txt",
						"ls -la",
					},
					Status: StatusSuccessful,
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "Step Name",
	      "image": "ubuntu:14.04",
	      "cmds": [
	        "ls -la",
	        "touch hello.txt",
	        "ls -la"
	      ]
			},
			{
	      "name": "Other Step Name",
	      "image": "ubuntu:14.04",
	      "cmds": [
	        "ls -la",
	        "touch hello.txt",
	        "ls -la"
	      ]
			}
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "Step Name",
					ImageName: "ubuntu:14.04",
					Cmds: []Cmd{
						"ls -la",
						"touch hello.txt",
						"ls -la",
					},
					Status: StatusSuccessful,
				},
				&Step{
					Name:      "Other Step Name",
					ImageName: "ubuntu:14.04",
					Cmds: []Cmd{
						"ls -la",
						"touch hello.txt",
						"ls -la",
					},
					Status: StatusSuccessful,
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step2"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step2"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1", "step2"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1", "step2"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
			},
		},
	},
	// Failure tests
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "Step Name",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls notafile.txt"]
			}
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "Step Name",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile.txt"},
					Status:    StatusFailed,
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls notafile"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile"},
					Status:    StatusFailed,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusNotRun,
					After:     []string{"step1"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls notafile"],
				"after": ["step1"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile"},
					Status:    StatusFailed,
					After:     []string{"step1"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls notafile"],
				"after": ["step1", "step2"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile"},
					Status:    StatusFailed,
					After:     []string{"step1", "step2"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["sleep 3"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls notafile"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1", "step2"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"sleep 3"},
					Status:    StatusStopped,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile"},
					Status:    StatusFailed,
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusNotRun,
					After:     []string{"step1", "step2"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls notafile"],
				"after": ["step1", "step2"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile"},
					Status:    StatusFailed,
					After:     []string{"step1", "step2"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
	  "name": "Pipeline Name",
	  "steps": [
	    {
	      "name": "step1",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls notafile"]
	    },
			{
	      "name": "step2",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    },
			{
	      "name": "step3",
	      "image": "ubuntu:14.04",
	      "cmds": ["ls"],
				"after": ["step1"]
	    }
	  ]
	}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile"},
					Status:    StatusFailed,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusNotRun,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusNotRun,
					After:     []string{"step1"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
		  "name": "Pipeline Name",
		  "steps": [
		    {
		      "name": "step1",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls"]
		    },
				{
		      "name": "step2",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls notafile"],
					"after": ["step1"]
		    },
				{
		      "name": "step3",
		      "image": "ubuntu:14.04",
		      "cmds": ["sleep 3"],
					"after": ["step1"]
		    }
		  ]
		}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile"},
					Status:    StatusFailed,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"sleep 3"},
					Status:    StatusStopped,
					After:     []string{"step1"},
				},
			},
		},
	},
	apiTestCase{
		requestBody: `{
		  "name": "Pipeline Name",
		  "steps": [
		    {
		      "name": "step1",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls"]
		    },
				{
		      "name": "step2",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls notafile"],
					"after": ["step1"]
		    },
				{
		      "name": "step3",
		      "image": "ubuntu:14.04",
		      "cmds": ["sleep 3"]
		    }
		  ]
		}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusFailed,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls notafile"},
					Status:    StatusFailed,
					After:     []string{"step1"},
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"sleep 3"},
					Status:    StatusStopped,
				},
			},
		},
	},
	// more complicated dependency tests
	apiTestCase{
		requestBody: `{
		  "name": "Pipeline Name",
		  "steps": [
		    {
		      "name": "step1",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls"]
		    },
				{
		      "name": "step2",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls"]
		    },
				{
		      "name": "step3",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls"],
					"after": ["step1", "step2"]
		    },
				{
		      "name": "step4",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls"],
					"after": ["step2"]
		    },
				{
		      "name": "step5",
		      "image": "ubuntu:14.04",
		      "cmds": ["ls"],
					"after": ["step3", "step4"]
		    }
		  ]
		}`,
		pipeline: Pipeline{
			Name:   "Pipeline Name",
			Status: StatusSuccessful,
			Steps: []*Step{
				&Step{
					Name:      "step1",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step2",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
				},
				&Step{
					Name:      "step3",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step1", "step2"},
				},
				&Step{
					Name:      "step4",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step2"},
				},
				&Step{
					Name:      "step5",
					ImageName: "ubuntu:14.04",
					Cmds:      []Cmd{"ls"},
					Status:    StatusSuccessful,
					After:     []string{"step3", "step4"},
				},
			},
		},
	},
}
