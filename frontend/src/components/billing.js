import React, {Component} from 'react';
import { useCookies } from 'react-cookie';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import "./billing.css";
import cfg from '../cfg';

export default function Billing() {

    const [cookies, setCookie] = useCookies(['jwt', 'jwtExpiration']);

    const oneDay = 24 * 60 * 60 * 1000; // hours*minutes*seconds*milliseconds
    const billingStartDate = "2022-03-21";
    const daysSinceBillingStart = Math.round(Math.abs((new Date().getTime() - new Date(billingStartDate).getTime()) / oneDay));

    class Billing extends Component {
        constructor() {
            super();
            this.state = {
                analytics: {
                    grandTotal: 0.00,
                    totalByUsers: [
                        /*
                        {
                            name: "example",
                            total: 829.00,
                            highestTotal: 110.00,
                            lastTotal: 0.00,
                            last_visit: "2022-03-20T15:43:01Z"
                        },
                        */
                    ],
                    userTotal: {},
                    userTotalsByTime: {
                        "last48hTotal": 0.0,
                        "lastWeekTotal": 0.0,
                        "thisWeekTotal": 0.0,
                        "thisMonthTotal": 0.0
                    }
                },
                showPopupAddBill: false,
                showInvalidAddBillResponse: false,
                inputAddBillAmount: ""
            };
            this.handleAddBillSubmit = this.handleAddBillSubmit.bind(this);
            this.handleAddBillInputChanged = this.handleAddBillInputChanged.bind(this);
            this.handleAddBillClick = this.handleAddBillClick.bind(this);
            this.handleAddBillClose = this.handleAddBillClose.bind(this);
            this.handleBillUndo = this.handleBillUndo.bind(this);
        }

        componentDidMount() {
            const requestOptions = {
                method: 'GET',
                headers: {'Authorization': 'Bearer ' + cookies.jwt}
            }
            fetch(cfg.BACKEND_HOST + '/bills', requestOptions)
                .then(res => res.json())
                .then(res => this.updateStateFromResponse(res));
        }
    
        updateStateFromResponse(res) {
            this.setState({
                analytics: res,
                showPopupAddBill: false,
                showInvalidAddBillResponse: false,
                inputAddBillAmount: "",
            });
        }
    
        handleAddBillSubmit(e) {
            e.preventDefault();
            const requestOptions = {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + cookies.jwt
                },
                body: JSON.stringify({
                    total: parseFloat(this.state.inputAddBillAmount.replace(",", "."))
                })
            };
            fetch(cfg.BACKEND_HOST + '/bills', requestOptions)
                .then(res => {
                    if (!res.ok) {
                        let newState = this.state;
                        newState.showInvalidAddBillResponse = true;
                        this.setState(newState);
                        throw new Error("Invalid request");
                    }   
                    return res.json();                 
                })
                .then(res => this.updateStateFromResponse(res))
                .catch(err => console.log(err));
        }

        handleAddBillInputChanged(e) {
            let newState  = this.state;
            newState.inputAddBillAmount = e.target.value;
            this.setState(newState);
        }

        handleAddBillClick(e) {
            let newState  = this.state;
            newState.showPopupAddBill = true;
            newState.showInvalidAddBillResponse = false;
            this.setState(newState);
        }

        handleAddBillClose(e) {
            let newState  = this.state;
            newState.showPopupAddBill = false;
            newState.showInvalidAddBillResponse = false;
            this.setState(newState);
        }

        handleBillUndo(e) {
            if (!window.confirm('Voulez-vous vraiment annuler la dernière facture? Le montant est de ' + this.state.analytics.userTotal.lastTotal.toFixed(2) + "$. Vous ne pouvez annuler qu'une seule facture.")) {
                return;
            }
            const requestOptions = {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + cookies.jwt
                },
                body: JSON.stringify()
            };
            fetch(cfg.BACKEND_HOST + '/bills/undo', requestOptions)
                .then(res => res.json())
                .then(res => this.updateStateFromResponse(res));
        }

        getFormattedDate(date) {
            if (date === "" || date == null) {
                return "";
            }
            let d = new Date(date);
            return d.toLocaleDateString("en-CA");
        }

        render() {
            let analyticsTableBody = this.state.analytics.totalByUsers.map((userTotal, index) =>  
                <tr key={index}>
                    <td >{userTotal.name}</td>
                    <td >{userTotal.total.toFixed(2)}$</td>
                    <td >{userTotal.highestTotal.toFixed(2)}$</td>
                    <td >{this.getFormattedDate(userTotal.lastVisit)}</td>
                </tr>
            );
    
            return (
                <div className="col-12">
                    <div className="fs-1 mt-5 mb-3 lh-sm">
                        <p className="mt-3 mb-3">🏴‍☠️🍻 Total Corsaire 🍻🏴‍☠️</p>
                        <h1 className="display-1">{this.state.analytics.grandTotal.toFixed(2)}$</h1>
                    </div>  
                    <div>
                        <p className="mb-0">Depuis {daysSinceBillingStart} jours</p>
                        <p>(21 mars 2022)</p>
                    </div>
                    <Button className="col-6 col-lg-3 m-0" onClick={this.handleAddBillClick} variant="dark">Ajout facture</Button>
                    
                    <>
                        <Modal show={this.state.showPopupAddBill} onHide={this.handleAddBillClose}>               
                            <Modal.Header closeButton>
                                <Modal.Title>Ajout facture</Modal.Title>
                            </Modal.Header>
                            <Modal.Body>
                                <form onSubmit={this.handleAddBillSubmit}>    
                                    <div className="row g-3 mb-3">
                                        <div className='col-auto'>
                                            <label htmlFor='inputBillAmount' className="col-form-label">Montant</label>
                                        </div>
                                        <div className='col-auto'>
                                            <input className="form-control" id="inputBillAmount" type="number" step="0.01" inputmode="decimal" aria-describedby="amount" placeholder="0.00" onChange={this.handleAddBillInputChanged}/>
                                        </div>
                                    </div>
                                    <h6 className={this.state.showInvalidAddBillResponse ? "text-danger mb-3" : "d-none"}>Requête invalide. Veuillez saisir un nombre valide!</h6>
                                    <div className="d-grid d-md-flex justify-content-center">
                                        <Button type="submit" onClick={this.handleAddBillSubmit} variant="dark">Ajouter</Button>

                                    </div>
                                    
                                </form>
                            </Modal.Body>
                            
                        </Modal>
                    </>
                            
                    
                    <div className="m-2">
                        <Button className="btn btn-sm" onClick={this.handleBillUndo} variant="outline-danger">Annuler</Button>
                    </div>

                    <div className="fs-1 mt-5 lh-sm">
                        <p className="mb-4">🍹 Dépenses Récentes 📈</p>
                        <div className='container'>
                            <div className='row d-flex justify-content-center'>
                                <div className="col-12 col-lg-5 mb-2">
                                    <div className="pt-1 border rounded">
                                        <h5>En 48h</h5>
                                        <p className='fs-2'>{this.state.analytics.userTotalsByTime.last48hTotal.toFixed(2)}$</p>
                                    </div>
                                </div>
                                <div className="col-12 col-lg-5 mb-2">
                                    <div className="pt-1 border rounded">
                                        <h5>Cette semaine</h5>
                                        <p className='fs-2'>{this.state.analytics.userTotalsByTime.thisWeekTotal.toFixed(2)}$</p>
                                    </div>
                                </div>
                                <div className="col-12 col-lg-5 mb-2">
                                    <div className="pt-1 border rounded">
                                        <h5>Ce mois-ci</h5>
                                        <p className='fs-2'>{this.state.analytics.userTotalsByTime.thisMonthTotal.toFixed(2)}$</p>
                                    </div>
                                </div>
                                <div className="col-12 col-lg-5 mb-2">
                                    <div className="pt-1 border rounded">
                                        <h5 >La semaine passée</h5>
                                        <p className='fs-2'>{this.state.analytics.userTotalsByTime.lastWeekTotal.toFixed(2)}$</p>
                                    </div>
                                </div>                  
                            </div>
                        </div>
                        
                    </div>  
                        
                    <div className="container mt-3">
                        <div className="row justify-content-center">
                            <div className='col-12 border border-2 rounded-3 align-self-center'>
                                <h2 className="p-3">Les légendes du Corsaire</h2>
                                <table className="table Billing-Table">
                                    <thead>
                                        <tr>
                                            <th scope='col'>Nom</th>
                                            <th scope='col'>Total</th>
                                            <th scope='col'>Record</th>
                                            <th scope='col'>Dernière visite</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {analyticsTableBody}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>

                </div>
            );
        }
    }
    return <Billing/>  
};