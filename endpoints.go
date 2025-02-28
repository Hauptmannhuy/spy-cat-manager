package main

import (
	"net/http"
)

func (app *application) getSpies(w http.ResponseWriter, r *http.Request) error {
	crudExec := crudExecutor[Spy]{
		handler: getAllHandler[Spy]{
			getAll: app.db.getSpies,
		},
	}
	json, err := crudExec.getAll()
	if err != nil {
		return err
	}
	sendResponse(w, 200, json)
	return nil
}

func (app *application) getSpy(w http.ResponseWriter, r *http.Request) error {

	crudExec := crudExecutor[Spy]{
		handler: getHandler[Spy]{
			fn: app.db.getSpy,
			r:  r,
			id: "id",
		},
	}
	responseData, err := crudExec.get()
	if err != nil {
		return err
	}
	sendResponse(w, http.StatusOK, responseData)
	return nil
}

func (app *application) createSpy(w http.ResponseWriter, r *http.Request) error {
	crudExec := crudExecutor[Spy]{
		handler: createHandler[Spy]{
			create: app.db.createSpy,
			r:      r,
		},
	}
	err := crudExec.create()

	if err != nil {
		return err
	}

	sendResponse(w, http.StatusOK)
	return nil
}

func (app *application) updateSpy(w http.ResponseWriter, r *http.Request) error {
	exec := crudExecutor[Spy]{
		handler: updateHandler[Spy]{
			update: app.db.updateSpy,
			r:      r,
			id:     "id",
		},
	}
	json, err := exec.update()
	if err != nil {
		return err
	}
	sendResponse(w, http.StatusOK, json)
	return nil
}

func (app *application) deleteSpy(w http.ResponseWriter, r *http.Request) error {
	crudExec := crudExecutor[any]{
		handler: deleteHandler[any]{
			delete: app.db.deleteSpy,
			r:      r,
			id:     "id",
		},
	}
	err := crudExec.delete()
	if err != nil {
		return err
	}
	sendResponse(w, http.StatusOK)
	return nil
}

func (app *application) getMission(w http.ResponseWriter, r *http.Request) error {
	crudExec := crudExecutor[Mission]{
		handler: getHandler[Mission]{
			fn: app.db.getMission,
			r:  r,
			id: "id",
		},
	}
	json, err := crudExec.get()
	if err != nil {
		return err
	}
	sendResponse(w, 200, json)
	return nil
}

func (app *application) getMissions(w http.ResponseWriter, r *http.Request) error {
	crudExec := crudExecutor[Mission]{
		handler: getAllHandler[Mission]{
			getAll: app.db.getMissions,
		},
	}
	json, err := crudExec.getAll()
	if err != nil {
		return err
	}
	sendResponse(w, 200, json)
	return nil
}

func (app *application) createMission(w http.ResponseWriter, r *http.Request) error {
	crudExec := crudExecutor[Mission]{
		handler: createHandler[Mission]{
			create: app.db.createMission,
			r:      r,
		},
	}
	err := crudExec.create()
	if err != nil {
		return err

	}
	sendResponse(w, 200)

	return nil
}

func (app *application) updateMission(w http.ResponseWriter, r *http.Request) error {
	m := &Mission{}
	err := decodeBody(r.Body, m)
	if err != nil {
		return err
	}

	if !m.checkUpdateValidity() {
		return serverError{
			reason: "mission state can only be set to true",
			code:   422,
		}
	}

	id := getURLvar(r, "id")
	err = app.db.updateMission(id, m)
	if err != nil {
		return dbError{
			reason: err.Error(),
			code:   400,
		}
	}

	sendResponse(w, 200)
	return nil
}

func (app *application) deleteMission(w http.ResponseWriter, r *http.Request) error {
	id := getURLvar(r, "id")
	err := app.db.deleteMission(id)
	if err != nil {
		return dbError{reason: err.Error(), code: 400}
	}
	sendResponse(w, 200)
	return nil
}

func (app *application) addTargetToMission(w http.ResponseWriter, r *http.Request) error {
	id := getURLvar(r, "id")
	crudExecutor := crudExecutor[Target]{
		handler: createHandler[Target]{
			r: r,
			create: func(data Target) error {
				return app.db.addTargetToMission(id, data)
			},
		},
	}
	err := crudExecutor.create()
	if err != nil {
		return err
	}
	sendResponse(w, 200)
	return nil
}

func (app *application) updateMissionTarget(w http.ResponseWriter, r *http.Request) error {
	mID := getURLvar(r, "mission_id")
	tID := getURLvar(r, "target_id")
	crudExecutor := crudExecutor[Target]{
		handler: updateHandler[Target]{
			r: r,
			update: func(i int, t Target) (Target, error) {
				return app.db.updateTarget(mID, tID, t)
			},
			id: "id",
		},
	}
	json, err := crudExecutor.update()
	if err != nil {
		return err
	}
	sendResponse(w, 200, json)
	return nil
}

func (app *application) deleteMissionTarget(w http.ResponseWriter, r *http.Request) error {
	mID := getURLvar(r, "mission_id")
	tID := getURLvar(r, "target_id")
	crudExecutor := crudExecutor[Target]{
		handler: deleteHandler[Target]{
			r: r,
			delete: func(i int) error {
				return app.db.deleteTarget(mID, tID)
			},
			id: "id",
		},
	}
	err := crudExecutor.delete()
	if err != nil {
		return err
	}
	sendResponse(w, 200)
	return nil
}
