import { TextField } from '@mui/material';
import React from 'react'
import { useFormContext, Controller } from 'react-hook-form'


const SubComponent1 = () => {
    const {
        control,
        formState: { errors },
    } = useFormContext();
    return (
        <>
            <Controller
                name="PublicKey"
                control={control}
                defaultValue=""
                render={({ field }) => (
                    <TextField
                        {...field}
                        label="Public Key"
                        variant="outlined"
                        error={!!errors.PublicKey}
                        helperText={
                            errors.PublicKey
                                ? errors.PublicKey?.message
                                : ''
                        }
                    />
                )}
            />
            <Controller
                name="PrivateKey"
                control={control}
                defaultValue=""
                render={({ field }) => (
                    <TextField
                        {...field}
                        label="Private Key"
                        variant="outlined"
                        error={!!errors.PrivateKey}
                        helperText={
                            errors.PrivateKey
                                ? errors.PrivateKey?.message
                                : ''
                        }
                    />
                )}
            />
            <Controller
                name="BlockchainAddress"
                control={control}
                defaultValue=""
                render={({ field }) => (
                    <TextField
                        {...field}
                        label="Blockchain Address"
                        variant="outlined"
                        error={!!errors.BlockchainAddress}
                        helperText={
                            errors.BlockchainAddress
                                ? errors.BlockchainAddress?.message
                                : ''
                        }
                    />
                )}
            />
            < input type="submit" />
        </>

    )
}
export default SubComponent1